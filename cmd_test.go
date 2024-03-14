package main

// import (
// 	"reflect"
// 	"testing"
// )

// func TestParser(t *testing.T) {
// 	// os.Args[1] = "aws"
// 	// os.Args[2] = "ec2"
// 	// os.Args[3] = "ec2"
// 	// os.Args[4] = "ec2"
// 	e, got := parser()
// 	want := []string{"aws", "ec2", "ec2"}

// 	if e != false {
// 		t.Errorf(
// 			"failed to return got: \nwant:%t\ngot:%t",
// 			false,
// 			e,
// 		)
// 	}
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf(
// 			"failed to return got: \nwant:%s\ngot:%s",
// 			want,
// 			got,
// 		)
// 	}
// }

// func TestParser1(t *testing.T) {

// 	os.Args[1] = "-e"
// 	os.Args[2] = "aws"
// 	os.Args[3] = "ec2"
// 	e, got := parser()
// 	want := []string{"aws", "ec2"}

// 	if e != true {
// 		t.Errorf(
// 			"failed to return got: \nwant:%t\ngot:%t",
// 			true,
// 			e,
// 		)
// 	}
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf(
// 			"failed to return got: \nwant:%s\ngot:%s",
// 			want,
// 			got,
// 		)
// 	}
// }
