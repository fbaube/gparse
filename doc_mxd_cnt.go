package gparse
/*
About Mixed Content

https://seanmcgrath.blogspot.com/2007/01/mixed-content-trying-to-understand-json.html

If you absolutely, totally, never, ever will need mixed content then there
are sane alternatives to XML. There always has been. from humble CSV up to
fancier JSON/Python/Ruby direct data expression languages.

A huge chunk of the world doesn't need mixed content or even know what it
is. They are the folks who look at the XML APIs and wonder "why is this so
difficult?"

XML is, and always was, a document-centric data representation language.

Hence mixed content. Hence XML. If you need mixed content you really need it.
If you don't need it, sometimes you cannot even conceptualise the problem it
solves. And yes, mixed content totally complicates the lives of those who
are using XML for data-centric applications.

The "standard" APIs of DOM and SAX handle the general case. They are extremely
sub-optimal for the very common data-centric case. We have no current standard
way to differentiate the former from the latter.

It would be a shame if this resulted in a "fork" in the road with fielded data,
yet again, going off on its own trajectory with document-centric data staying
on the XML road. Too much good stuff to be lost that way.

The nut that needs to be cracked to stop this happening is the Mixed Content
Case nut.

[ Go takes the fielded data fork. ]
*/
