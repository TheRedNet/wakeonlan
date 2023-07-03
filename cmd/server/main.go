package main

//he
import (
	"fmt"
	"net"

	"github.com/TheRedNet/wakeonlan/embedFS"
	"github.com/TheRedNet/wakeonlan/pkg/config"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	wol "github.com/sabhiram/go-wol/wol"
)

//go:generate go run github.com/UnnoTed/fileb0x b0x.yaml
var editMode bool
var lastError string

func sendMagicPacket(macAddr, bcastAddr string) (err error) {
	packet, err := wol.New(macAddr)
	if err != nil {
		return
	}
	payload, err := packet.Marshal()
	if err != nil {
		return
	}

	if bcastAddr == "" {
		bcastAddr = "255.255.255.255:9"
	}
	bcAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, bcAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Write(payload)
	return
}

func handleToggle(ctx *fiber.Ctx) error {
	editMode = !editMode
	deviceMap, err := config.LoadDevices()
	if err != nil {
		lastError = fmt.Sprintf("Failed to load devices.json: %v", err)
		return ctx.Redirect("/")
	}
	var devices []config.Device
	for name, mac := range deviceMap {
		devices = append(devices, config.Device{Name: name, Mac: mac})
	}
	return ctx.Redirect("/")
}

func handleEdit(ctx *fiber.Ctx) error {
	action := ctx.FormValue("action")
	name := ctx.FormValue("name")
	mac := ctx.FormValue("mac")
	deviceMap, err := config.LoadDevices()
	if err != nil {
		deviceMap = make(map[string]string)
		lastError = fmt.Sprintf("Failed to load devices.json: %v. Created a new one", err)
		//return ctx.Redirect("/")
	}
	switch action {
	case "add":
		deviceMap[name] = mac
	case "delete":
		delete(deviceMap, name)
	default:
		lastError = fmt.Sprintf("Invalid action: %s", action)
		return ctx.Redirect("/")
	}
	err = config.SaveDevices(deviceMap)
	if err != nil {
		lastError = fmt.Sprintf("Failed to save devices.json: %v", err)
		return ctx.Redirect("/")
	}
	return ctx.Redirect("/")
}

func handleWake(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name")
	deviceMap, err := config.LoadDevices()
	if err != nil {
		lastError = fmt.Sprintf("Failed to load devices.json: %v", err)
		return ctx.Redirect("/")
	}
	mac := deviceMap[name]
	if mac == "" {
		lastError = fmt.Sprintf("Unknown device: %s", name)
		return ctx.Redirect("/")
	}
	err = sendMagicPacket(mac, "255.255.255.255:9")
	if err != nil {
		lastError = fmt.Sprintf("Failed to send magic packet: %v", err)
		return ctx.Redirect("/")
	}
	return ctx.Redirect("/")
}

// handles calls to "/" and renders the index page.
func handleIndex(ctx *fiber.Ctx) error {
	deviceMap, err := config.LoadDevices()
	if err != nil {
		lastError = fmt.Sprintf("Failed to load config.json: %v", err)
		//return ctx.Redirect("/")
	}
	var devices []config.Device
	for name, mac := range deviceMap {
		devices = append(devices, config.Device{Name: name, Mac: mac})
	}
	displayedError := lastError
	lastError = ""
	hasErrored := displayedError != ""
	editModebool := editMode
	if hasErrored {
		println(displayedError)
	}
	return ctx.Render("views/index", fiber.Map{
		"Devices":    devices,
		"EditMode":   editModebool,
		"LastError":  displayedError,
		"HasErrored": hasErrored,
	})
}

func main() {
	fmt.Println("---- Preparing server ----")
	fmt.Println("Creating engine...")
	engine := html.NewFileSystem(embedFS.HTTP, ".html")
	fmt.Println("Creating app...")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	fmt.Println("Setting up routes...")
	fmt.Println("Route: GET /")
	app.Get("/", handleIndex)
	fmt.Println("Route: POST /wake")
	app.Post("/wake", handleWake)
	fmt.Println("Route: POST /toggle")
	app.Post("/toggle", handleToggle)
	fmt.Println("Route: POST /edit")
	app.Post("/edit", handleEdit)
	fmt.Println("---- Starting server on port 8000 ----")
	err := app.Listen("0.0.0.0:8000")
	if err != nil {
		fmt.Printf("Server failed: %v/n", err)
	}
	fmt.Println("---- Server stopped ----")
}
