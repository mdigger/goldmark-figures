package figures_test

import (
	"log"
	"os"

	figures "github.com/mdigger/goldmark-figures"
	"github.com/yuin/goldmark"
)

func Example() {
	var source = []byte(`
![**Figure:** [description](/link)](image.png "title")

![](image.png "title")

![**Figure:** [description](/link)](image.png "title")
![**alt**](image.png "title")
`)
	var md = goldmark.New(figures.Enable)
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// <figure>
	// <img src="image.png" alt="Figure: description" title="title">
	// <figcaption><strong>Figure:</strong> <a href="/link">description</a></figcaption>
	// </figure>
	// <p><img src="image.png" alt="" title="title"></p>
	// <p><img src="image.png" alt="Figure: description" title="title">
	// <img src="image.png" alt="alt" title="title"></p>
}
