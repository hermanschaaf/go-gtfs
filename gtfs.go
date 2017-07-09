package gtfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/artonge/go-csv-tag"
)

// Load - load GTFS files
// @param dirPath: the directory containing the GTFS
// @return a filled GTFS or an error
func Load(dirPath string) (*GTFS, error) {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Error loading GTFS: directory does not existe")
	}
	g := &GTFS{Path: dirPath}
	err = loadGTFS(g)
	if err != nil {
		return nil, fmt.Errorf("Error loading GTFS: '%v'\n	==> %v", g.Path, err)
	}
	return g, nil
}

// LoadSplitted - load splitted GTFS files
// ==> When GTFS are splitted into sub directories
// @param dirPath: the directory containing the sub GTFSs
// @return an array of filled GTFS or an error
func LoadSplitted(dirPath string) ([]*GTFS, error) {
	// Get directory list
	subDirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	// Prepare the array of GTFSs
	GTFSs := make([]*GTFS, len(subDirs))
	i := 0
	// Load each sub directory into a new GTFS struct
	for _, dir := range subDirs {
		if !dir.IsDir() {
			continue
		}
		GTFSs[i] = &GTFS{Path: path.Join(dirPath, dir.Name())}
		err := loadGTFS(GTFSs[i])
		if err != nil {
			return nil, fmt.Errorf("Error loading GTFS: '%v'\n	==> %v", GTFSs[i].Path, err)
		}
		i++
	}
	// Return a slice of GTFSs, because we skipped non directory files
	return GTFSs[:i], nil
}

// Load a directory containing:
// 	- routes.txt
// 	- stops.txt
// @param g: the GTFS struct that will receive the data
// @return an error
func loadGTFS(g *GTFS) error {
	// Create a slice of agency to load agency.txt
	var agencySlice []Agency
	// List all files that will be loaded and there dest
	filesToLoad := map[string]interface{}{
		"agency.txt":     &agencySlice,
		"calendar.txt":   &g.Calendars,
		"routes.txt":     &g.Routes,
		"stops.txt":      &g.Stops,
		"stop_times.txt": &g.StopsTimes,
		"transfers.txt":  &g.Transfers,
		"trips.txt":      &g.Trips,
	}
	// Load the files
	for file, dest := range filesToLoad {
		filePath := path.Join(g.Path, file)
		// If the file does not existe, skip it
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			continue
		}
		err = csvtag.Load(csvtag.Config{
			Path: filePath,
			Dest: dest,
		})
		if err != nil {
			return fmt.Errorf("Error loading file (%v)\n	==> %v", file, err)
		}
	}
	// Put the loaded agency in g.Agency
	if len(agencySlice) > 0 {
		g.Agency = agencySlice[0]
	}
	return nil
}