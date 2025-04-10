package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	System32_files_to_dump        []string
	System32Drivers_files_to_dump []string
	Silent                        bool
}

var config Config

func checkIfFilepathExists(file_path string) bool {
	_, err := os.Stat(file_path)
	return !os.IsNotExist(err)
}

func getConfig() (Config, error) {
	config_raw_data, err := os.ReadFile("config.toml")
	if err != nil {
		return Config{}, err
	}

	config_string := string(config_raw_data)

	var config Config
	err = toml.Unmarshal([]byte(config_string), &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func ShouldDumpAll(selected_files []string) bool {
	files_amount := len(selected_files)
	if files_amount == 1 && selected_files[0] == "*" {
		return true
	}
	return false
}

func createFolder(dir_name string) (string, error) {
	err := os.Mkdir(dir_name, 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return dir_name, nil
}

func copyFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func hasProperExtension(file os.DirEntry) bool {
	return filepath.Ext(file.Name()) == ".sys" || filepath.Ext(file.Name()) == ".exe" || filepath.Ext(file.Name()) == ".dll"
}

func dumpEveryFileInDirectory(source_directory string, destination_directory string) {
	files, err := os.ReadDir(source_directory)
	if err != nil {
		fmt.Println("[error] Couldn't retrieve files from", source_directory)
		return
	}

	copied_counter := 0
	copyable_files := 0
	for _, file := range files {
		if hasProperExtension(file) {
			src := filepath.Join(source_directory, file.Name())
			dest := filepath.Join(destination_directory, file.Name())

			copy_err := copyFile(src, dest)
			if copy_err != nil {
				fmt.Println("[error] Couldn't copy:", src, "Error:", copy_err)
			} else {
				if !config.Silent {
					fmt.Println("[+]", src, "copied.")
				}
				copied_counter++
			}
			copyable_files++
		}
	}

	fmt.Printf("[+] Copied %d/%d files from %s\n", copied_counter, copyable_files, source_directory)
}

func dumpSelectedFilesInDirectory(source_directory string, destination_directory string, selected_files []string) {
	copied_counter := 0
	for _, v := range selected_files {
		file_path := filepath.Join(source_directory, v)
		if checkIfFilepathExists(file_path) {
			copy_err := copyFile(file_path, filepath.Join(destination_directory, v))
			if copy_err != nil {
				fmt.Println("[error] Couldn't copy:", file_path, "Error:", copy_err)
			} else {
				if !config.Silent {
					fmt.Println("[+]", file_path, "copied.")
				}
				copied_counter++
			}
		} else {
			fmt.Println("[warning]", file_path, "does not exist.")
		}
	}

	fmt.Printf("[+] Copied %d/%d files from %s\n", copied_counter, len(selected_files), source_directory)
}

func dumpDirectory(source_directory string, destination_directory string, selected_modules []string) {
	if ShouldDumpAll(selected_modules) {
		dumpEveryFileInDirectory(source_directory, destination_directory)
	} else {
		dumpSelectedFilesInDirectory(source_directory, destination_directory, selected_modules)
	}
}

func main() {
	if !checkIfFilepathExists("config.toml") {
		fmt.Println("[error] Couldn't find config.toml")
		return
	}

	var err error
	config, err = getConfig()
	if err != nil {
		fmt.Println("[error] Confing parsing failed. Error: ", err)
		return
	}

	dumps_dir, err := createFolder("dumps")
	if err != nil {
		fmt.Println("[error] Couldn't create dumps folder. Error: ", err)
		return
	}

	current_date_as_string := time.Now().Local().Format("02-01-2006_15-04-05")

	current_dump_dir, err := createFolder(filepath.Join(dumps_dir, current_date_as_string))
	if err != nil {
		fmt.Printf("[error] Couldn't create dumps/%s/ folder. Error: %s\n", current_date_as_string, err)
		return
	}

	current_dump_drivers_dir, err := createFolder(filepath.Join(current_dump_dir, "drivers"))
	if err != nil {
		fmt.Printf("[error] Couldn't create dumps/%s/drivers folder. Error: %s\n", current_date_as_string, err)
		return
	}

	system_root_path := strings.ToLower(os.Getenv("systemroot"))
	if system_root_path == "" {
		fmt.Println("[error] Couldn't retrieve system root directory")
		return
	}

	system32_path := filepath.Join(system_root_path, "System32")
	dumpDirectory(system32_path, current_dump_dir, config.System32_files_to_dump)

	system32_drivers_path := filepath.Join(system32_path, "drivers")
	dumpDirectory(system32_drivers_path, current_dump_drivers_dir, config.System32Drivers_files_to_dump)

	fmt.Println("[+] Dumping process has been finished")
}
