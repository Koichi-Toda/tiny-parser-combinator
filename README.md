# tiny-parser-combinator
This is very small parser combinator sample codes for giving essential concept.

# Basic Sample
You can define parser combinator as following. 

	expr = Seq(&term,
		Rep(Choice(
			Seq(e("+"), &term),
			Seq(e("-"), &term))))

	term = Seq(&factor,
		Rep(Choice(
			Seq(e("*"), &factor),
			Seq(e("/"), &factor))))
	factor = Choice(
		number,
		Seq(e("("), Seq(&expr, e(")"))))

You can run a generated parser using ParseAll function.

	parser.ParseAll(expr, "(1+2+3)*(4+5+6)")
  

And also you can transform parse result to get something you want.

	expr = T_(Seq(&term,
		Rep(Choice(
			T_(Seq(e("+"), &term), f1),
			T_(Seq(e("-"), &term), f1)))), f2)

	term = T_(Seq(&factor,
		Rep(Choice(
			T_(Seq(e("*"), &factor), f1),
			T_(Seq(e("/"), &factor), f1)))), f2)
	factor = Choice(
		number,
		T_(Seq(e("("), Seq(&expr, e(")"))), f3))
    
You can get calculate result.

	parser.ParseAll(expr, "(1+2+3)*(4+5+6)")
	
	$ go run main.go
	$ parsed: [Parse Success] accepted:(1+2+3)*(4+5+6),result:[90], rest:
  
More information can be found on the following paths.

https://github.com/Koichi-Toda/tiny-parser-combinator/blob/main/parser/example.go

