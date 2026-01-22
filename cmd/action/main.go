package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	common_entity "github.com/desain-gratis/common/types/entity"
	"github.com/desain-gratis/deployd/src/entity"
	"github.com/rs/zerolog/log"

	contentsync "github.com/desain-gratis/common/delivery/mycontent-api-client"
	"github.com/desain-gratis/deployd-gh-action/internal/utility"
)

func env(key string, required bool) string {
	v := os.Getenv(key)
	if v == "" && required {
		fmt.Fprintf(os.Stderr, "missing env %s\n", key)
		os.Exit(1)
	}
	return v
}

func main() {
	urlx := env("INPUT_URL", true)
	namespace := env("INPUT_NAMESPACE", true)
	name := env("INPUT_NAME", true)
	archive := env("INPUT_ARCHIVE", true)

	eventPath := env("GITHUB_EVENT_PATH", true)
	eventData, err := os.ReadFile(eventPath)
	if err != nil {
		log.Warn().Msgf("unable to read github event path: %v", err)
	}

	commitID := env("GITHUB_SHA", true)
	actor := env("GITHUB_ACTOR", true)

	ctx := context.Background()

	wd, err := os.Getwd()
	log.Info().Msgf("pwd: %v %v", wd, err)
	log.Info().Msgf("input name: %v", name)

	isDir, err := utility.IsDir(archive)
	if err != nil {
		log.Panic().Msgf("error reading archive: %v", err)
	}

	if isDir {
		outputArchive := "./tmp/archive.tgz"
		log.Info().Msgf("Bundling to tgz")
		err := utility.BundleDir(outputArchive, archive)
		if err != nil {
			log.Panic().Msgf("error bundling archive: %v", err)
		}
		archive = outputArchive
	}

	osArch := "linux/amd64"
	tags := []string{"os/arch:linux/amd64"}

	var branch string
	var tag string

	if env("GITHUB_REF_TYPE", true) == "branch" {
		branch = env("GITHUB_REF_NAME", true)
		tags = append(tags, fmt.Sprintf("branch:%v", branch))
	} else if env("GITHUB_REF_TYPE", true) == "tag" {
		tag = env("GITHUB_REF_NAME", true)
		tags = append(tags, fmt.Sprintf("tag:%v", tag))
	}

	u, err := url.Parse(urlx + "/artifactd/build")
	if err != nil {
		log.Panic().Msgf("error parsing url: %v", err)
	}

	// tagsStr := strings.Join(tags, ",")
	data := []*entity.Artifact{{
		UID:          commitID,
		Ns:           namespace,
		CommitID:     commitID,
		Branch:       branch,
		Actor:        actor,
		Tag:          tag,
		Data:         json.RawMessage(eventData),
		PublishedAt:  time.Now(),
		Source:       "github",
		RepositoryID: name,
		OsArch:       []string{osArch}, // hardcode first
		URLx:         "",
		Name:         name,
		Archive: []*common_entity.File{
			{Id: commitID + "|" + osArch, Url: archive},
		},
	}}

	buildSync := contentsync.Builder[*entity.Artifact](u, "repository").
		WithNamespace("*").
		WithData(data)

	// TODO: improve DevX it's a bit painful
	// because you can get parameter automatically from parents
	//

	// notice there is no "repository" here because the main entity Artifact already have "repository" ref.
	buildSync.
		WithFiles(getArchive, "../archive", "build")

	// upload metadata

	err = buildSync.Build().Execute(ctx)
	if err != nil {
		log.Panic().Msgf("failed to execute: %v", err)
	}

	fmt.Println("Archive uploaded successfully")
}

func getArchive(t []*entity.Artifact) []contentsync.FileContext[*entity.Artifact] {
	result := make([]contentsync.FileContext[*entity.Artifact], 0)
	for i := range t {
		for j := range t[i].Archive {
			result = append(result, contentsync.FileContext[*entity.Artifact]{
				Base: t[i], File: &t[i].Archive[j],
			})
		}
	}
	return result
}
