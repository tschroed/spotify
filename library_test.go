package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	kAlbumID0 = "4iV5W9uYEdYUVa79Axb7Rh"
	kAlbumID1 = "1301WleyT98MSxVHPZCA6M"
	kAlbumID2 = "0udZHhCi7p1YzMlvI4fXoK"
	kAlbumID3 = "55nlbqqFVnSsArIeYSQlqx"
)

func albumBodyValidator(t *testing.T) func(*http.Request) {
	return func(req *http.Request) {
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal("Could not read request body:", err)
		}

		var body map[string]interface{}
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			t.Fatal("Error decoding request body:", err)
		}
		idsArray, ok := body["ids"]
		if !ok {
			t.Error("No album IDs in request body")
		}
		idsSlice := idsArray.([]interface{})
		if l := len(idsSlice); l != 2 {
			t.Fatalf("Expected 2 albums, got %d\n", l)
		}
		id0 := idsSlice[0].(string)
		if id0 != kAlbumID0 {
			t.Errorf("[id0] Expected %s, got %s", kAlbumID0, id0)
		}
		id1 := idsSlice[1].(string)
		if id1 != kAlbumID1 {
			t.Errorf("[id1] Expected %s, got %s", kAlbumID1, id1)
		}
	}
}

func TestUserHasTracks(t *testing.T) {
	client, server := testClientString(http.StatusOK, `[ false, true ]`)
	defer server.Close()

	contains, err := client.UserHasTracks(context.Background(), kAlbumID2, kAlbumID3)
	if err != nil {
		t.Error(err)
	}
	if l := len(contains); l != 2 {
		t.Error("Expected 2 results, got", l)
	}
	if contains[0] || !contains[1] {
		t.Error("Expected [false, true], got", contains)
	}
}

func TestAddTracksToLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "", albumBodyValidator(t))
	defer server.Close()

	err := client.AddTracksToLibrary(context.Background(), kAlbumID0, kAlbumID1)
	if err != nil {
		t.Error(err)
	}
}

func TestAddTracksToLibraryFailure(t *testing.T) {
	client, server := testClientString(http.StatusUnauthorized, `
{
  "error": {
    "status": 401,
    "message": "Invalid access token"
  }
}`)
	defer server.Close()
	err := client.AddTracksToLibrary(context.Background(), kAlbumID0, kAlbumID1)
	if err == nil {
		t.Error("Expected error and didn't get one")
	}
}

func TestAddTracksToLibraryWithContextCancelled(t *testing.T) {
	client, server := testClientString(http.StatusOK, ``)
	defer server.Close()

	ctx, done := context.WithCancel(context.Background())
	done()

	err := client.AddTracksToLibrary(ctx, kAlbumID0, kAlbumID1)
	if !errors.Is(err, context.Canceled) {
		t.Error("Expected error and didn't get one")
	}
}

func TestRemoveTracksFromLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "", albumBodyValidator(t))
	defer server.Close()

	err := client.RemoveTracksFromLibrary(context.Background(), kAlbumID0, kAlbumID1)
	if err != nil {
		t.Error(err)
	}
}

func TestUserHasAlbums(t *testing.T) {
	client, server := testClientString(http.StatusOK, `[ false, true ]`)
	defer server.Close()

	contains, err := client.UserHasAlbums(context.Background(), kAlbumID2, kAlbumID3)
	if err != nil {
		t.Error(err)
	}
	if l := len(contains); l != 2 {
		t.Error("Expected 2 results, got", l)
	}
	if contains[0] || !contains[1] {
		t.Error("Expected [false, true], got", contains)
	}
}

func TestAddAlbumsToLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	err := client.AddAlbumsToLibrary(context.Background(), kAlbumID0, kAlbumID1)
	if err != nil {
		t.Error(err)
	}
}

func TestAddAlbumsToLibraryFailure(t *testing.T) {
	client, server := testClientString(http.StatusUnauthorized, `
{
  "error": {
    "status": 401,
    "message": "Invalid access token"
  }
}`)
	defer server.Close()
	err := client.AddAlbumsToLibrary(context.Background(), kAlbumID0, kAlbumID1)
	if err == nil {
		t.Error("Expected error and didn't get one")
	}
}

func TestRemoveAlbumsFromLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	err := client.RemoveAlbumsFromLibrary(context.Background(), kAlbumID0, kAlbumID1)
	if err != nil {
		t.Error(err)
	}
}
