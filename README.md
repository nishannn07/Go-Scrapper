# 🕷️ Go Web Scraper

A simple command-line web scraper written in Go, allowing you to extract links and headlines from any webpage. It uses the powerful [goquery](https://github.com/PuerkitoBio/goquery) library for HTML parsing.

---

## 📦 Features

- Extract all anchor links (`<a href="...">`)
- Extract all headlines (`<h1>`, `<h2>`, `<h3>`)
- Output to terminal or save results to a file
- Easy to use with command-line flags

---

## 🚀 Getting Started

### 1. Clone this repository

```bash
git clone https://github.com/your-username/go-web-scraper.git
cd go-web-scraper
```

### 2. Install dependencies

```bash
go get -u github.com/PuerkitoBio/goquery
```

---

## 🧪 Usage

```bash
go run main.go -url <target_url> [-extract links|headlines|all] [-output output.txt]
```

### Arguments:

| Flag         | Description                                                                 | Required | Example                                 |
|--------------|-----------------------------------------------------------------------------|----------|-----------------------------------------|
| `-url`       | Target URL to scrape                                                        | ✅ Yes   | `-url https://example.com`              |
| `-extract`   | Data type to extract: `links`, `headlines`, or `all` (default is `links`)   | ❌ No    | `-extract headlines`                    |
| `-output`    | File path to save the output (default is to print to console)               | ❌ No    | `-output results.txt`                   |

---

## 📌 Examples

### Extract links and print to console:
```bash
go run main.go -url https://example.com
```

### Extract headlines and save to a file:
```bash
go run main.go -url https://example.com -extract headlines -output headlines.txt
```

### Extract both links and headlines:
```bash
go run main.go -url https://example.com -extract all
```

---

## ⚠️ Disclaimer

> ❗ For educational purposes only.
>
> Do **NOT** scrape websites like **Amazon**, **Flipkart**, or any site that restricts automated access through their [Terms of Service](https://www.amazon.in/gp/help/customer/display.html?nodeId=201909000) or `robots.txt`.  
> Always get permission when scraping websites, and avoid overwhelming servers.

---

## 🧑‍💻 Author

Made with 💻 by Nishan

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
