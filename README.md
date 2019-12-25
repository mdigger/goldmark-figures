# goldmark-figures

This [goldmark](https://github.com/yuin/goldmark) parser extension adds 
paragraph image render as figure.

An image with nonempty alt text, occurring by itself in a paragraph, will be 
rendered as a figure with a caption. The image’s alt text will be used as the 
caption.

```markdown
![**Figure:** [description](/link)](image.png "title")
```

```html
<figure>
<img src="image.png" alt="Figure: description" title="title">
<figcaption><strong>Figure:</strong> <a href="/link">description</a></figcaption>
</figure>
```

If you just want a regular inline image, just make sure it is not the only thing 
in the paragraph. One way to do this is to insert a nonbreaking space after the 
image:

```markdown
![This image won't be a figure](/url/of/image.png)\
```

This syntax is borrowed from [Pandoc](https://pandoc.org/MANUAL.html#images).
