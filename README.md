# go-captcha
> 2captcha.com library written in go

## Installation
    go get github.com/dank/go-captcha

## Example
```go
c := captcha.New("youd9kodxeuwxn0gzvaancmrkt895ua9")

solved, err := c.Solve("https://www.google.com/recaptcha/api2/demo?invisible=true", "6LfP0CITAAAAAHq9FOgCo7v_fb0-pmmH9VW3ziFs", true)
if err != nil {
	log.Fatalln(err)
}

fmt.Println(solved)
```

TODO