package foldup

import gstorage "cloud.google.com/go/storage"

func revertStubs() {
	newGCSClient = gstorage.NewClient
}
