package main

func main() {
	a := App{}
	a.Initialize("root", "", "db_sim")
	a.Run(":8080")
}
