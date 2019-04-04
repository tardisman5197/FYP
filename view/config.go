package view

import "os"

// pathToDir is the directory path to the executable of the project
var pathToDir, _ = os.Getwd()

// pathToUnity is the file path to the unity application
var pathToUnity = pathToDir + "\\unity\\FYP.exe"

// pathToImages is the file path to the images that are created by unity
var pathToImages = pathToDir + "\\unity\\FYP_Data\\"

// removeImagesOnShutdown is true if the images created by unity
// should be removed when the server is shutdown
const removeImagesOnShutdown = true

// startUnity is a bool that is true if a unity application should be
// started when creating a unityServer
const startUnity = true
