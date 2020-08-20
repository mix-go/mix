package dotenv

import "github.com/joho/godotenv"

// Load will read your env file(s) and load them into ENV for this process.
func Load(filenames ...string) (err error) {
    return godotenv.Load(filenames...)
}
