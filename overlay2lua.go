package main

import (
	"bytes"
	"flag"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"image"
	icolor "image/color"
	png "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	json_data *simplejson.Json
}

type OffsetTable map[string][][]image.Point

func (config *Config) getTileGeometry(geom_type string) image.Point {
	tile_size := config.json_data.Get("tile_geometry").Get(geom_type)
	tile_size_x, err := tile_size.GetIndex(0).Int()
	if err != nil {
		log.Fatal(err)
	}
	tile_size_y, err := tile_size.GetIndex(1).Int()
	if err != nil {
		log.Fatal(err)
	}
	return image.Point{tile_size_x, tile_size_y}
}

func (config *Config) getColorMappings() map[string]string {
	json_map := config.json_data.Get("colors")
	raw_map, err := json_map.Map()
	if err != nil {
		log.Fatal(err)
	}

	str_map := make(map[string]string)
	for key, _ := range raw_map {
		str_map[key], err = json_map.Get(key).String()
		if err != nil {
			log.Fatal(err)
		}
	}
	return str_map
}

func (config *Config) getBlankTable(gridsize image.Point) OffsetTable {
	table_data := make(OffsetTable)
	for _, name := range config.getColorMappings() {
		array2D := make([][]image.Point, gridsize.Y)
		for i := 0; i < gridsize.Y; i++ {
			array2D[i] = make([]image.Point, gridsize.X)
		}
		table_data[name] = array2D
	}
	return table_data
}

func hexToColor(hex string) icolor.RGBA {
	rgb, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		log.Fatal(err)
	}

	r := uint8(rgb >> 16)
	g := uint8(rgb>>8) & 0xFF
	b := uint8(rgb) & 0xFF
	return icolor.RGBA{r, g, b, 0xFF}
}

func exportMatrices(filename string, config *Config) {
	fmt.Printf("Processing %s... ", filename)

	// Open image
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode image
	image_file, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Set up some image variables
	image_size := image_file.Bounds().Size()
	tile_size := config.getTileGeometry("size")
	tile_origin := config.getTileGeometry("origin")

	tiles_x := image_size.X / tile_size.X
	tiles_y := image_size.Y / tile_size.Y
	fmt.Printf("Sprite sheet size: %dx%d\n", tiles_x, tiles_y)

	// Set up some path variables
	luaFilename := strings.Replace(filepath.Base(filename), filepath.Ext(filename), ".lua", 1)
	outputPath := filepath.Join(filepath.Dir(filename), luaFilename)

	// Open output lua file
	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	// Populate table_data with blank data at the appropriate size
	table_data := config.getBlankTable(image.Point{tiles_x, tiles_y})

	// Populate color data for speed
	colors := make(map[icolor.RGBA]string)
	for hex, name := range config.getColorMappings() {
		color := hexToColor(hex)
		colors[color] = name
	}

	// Loop through the image
	for y := 0; y < image_size.Y; y++ {
		for x := 0; x < image_size.X; x++ {
			r, g, b, a := image_file.At(x, y).RGBA()
			color := icolor.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			name, ok := colors[color]
			if ok {
				current_tile := image.Point{
					x / tile_size.X,
					y / tile_size.Y,
				}
				current_offset := image.Point{
					x % tile_size.X,
					y % tile_size.Y,
				}
				table_data[name][current_tile.Y][current_tile.X] = current_offset.Sub(tile_origin)
			}
		}
	}

	// Write this data to disk
	for name, _ := range table_data {
		n, err := fmt.Fprintf(output, "%s = {", name) // start table
		if n == 0 || err != nil {
			log.Fatal(err)
		}

		for y := 0; y < tiles_y; y++ {
			fmt.Fprintf(output, " {") // start row
			for x := 0; x < tiles_x; x++ {
				point := table_data[name][y][x]
				fmt.Fprintf(output, " {%d, %d},", point.X, point.Y)
			}
			fmt.Fprintln(output, " },") // end row
		}

		fmt.Fprintln(output, " }\n") // End table
	}
}

func getConfig(filename string) (*Config, error) {
	// Load config file
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Read all file contents into a bytes.Buffer
	var contents bytes.Buffer
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}

		contents.Write(buf[:n])
	}

	json, err := simplejson.NewJson(contents.Bytes())
	config := new(Config)
	config.json_data = json
	return config, err
}

func main() {
	flag.Parse()

	config, err := getConfig(flag.Args()[0])
	if err == nil {
		fmt.Println("Config loaded successfully")
	} else {
		fmt.Println("Config not loaded successfully")
		fmt.Println(err)
		return
	}
	for _, handle := range flag.Args()[1:] {
		exportMatrices(handle, config)
	}
}
