package controllers

import "time"

// Now is a variable that points to time.Now, allowing override in tests
var Now = time.Now
