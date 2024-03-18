package main

// import (
// 	"testing"
// )

// func TestCmdInit(t *testing.T) {
// 	err := cmdInit("", "")
// 	if err == nil {
// 		t.Error(err)
// 	}
// 	err = cmdInit("", "yes")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestSelectShowType(t *testing.T) {
// 	stat := IsAFile
// 	err := selectShowType(stat)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// // func TestParser1(t *testing.T) {

// // 	os.Args[1] = "-e"
// // 	os.Args[2] = "aws"
// // 	os.Args[3] = "ec2"
// // 	e, got := parser()
// // 	want := []string{"aws", "ec2"}

// // 	if e != true {
// // 		t.Errorf(
// // 			"failed to return got: \nwant:%t\ngot:%t",
// // 			true,
// // 			e,
// // 		)
// // 	}
// // 	if !reflect.DeepEqual(got, want) {
// // 		t.Errorf(
// // 			"failed to return got: \nwant:%s\ngot:%s",
// // 			want,
// // 			got,
// // 		)
// // 	}
// // }
