package cli (
  "os"
)

type Args struct {
  help Help
  load Load
  last Last
}

func Args(a ...string) Args {
  f := args{}
  for i, arg := range a {
    switch {
    case arg == "help" || arg == "-h" || arg == "--help":
      f.help = Help{true}
    case args == "load" || args == "-l" || args == "--load":
      f.load == Load{a[i + 1]}
    case args == "last" || args == "-L" || args == "--last":
      f.last == Last{true}
    }
  }
}

type flag interface{
  parse() 
}

type Help bool{}

func (h Help) parse() {
  if h == true {
    println("
  }
}

type Last bool{}

type Load string{}
