package main

import (
  "fmt"
  "github.com/ywwg/packinglib"
)

func main() {
  t := packinglib.Trip{
    Days: 4,
    C: packinglib.FireflyContext,
  }

  for _, p := range t.MakeList() {
    if p.Units == packinglib.NoUnits {
      if p.Count == float64(int(p.Count)) {
        fmt.Printf("%d %s\n", int(p.Count), p.Name)
      } else {
        fmt.Printf("%.1f %s\n", p.Count, p.Name)
      }
    } else {
      if p.Count == float64(int(p.Count)) {
        fmt.Printf("%d %s of %s\n", int(p.Count), p.Units, p.Name)
      } else {
        fmt.Printf("%.1f %s of %s\n", p.Count, p.Units, p.Name)
      }
    }
  }
}
