package deployment

import (
	"os"
	"testing"

	"github.com/biensupernice/krane/internal/data"
)

// To signal failure in a test use `t.Error(), t.Errorf() t.Fail()`

const (
	templateWithNoImageAndName = `{}`

	templateWithNoImage = `
		{
			"name": "testing",
			"config": {}
		}
	`

	templateWithNoName = `
		{
			"config': {
				"image": "bsn/image"
			}
		}
	`
	templateWithNameAndImage = `
	{
		"name": "testing",
		"config": {	
			"image": "bsn/image"
		}
	}
	`
)

// Complete template to use for testing
var testTemplate = &Template{
	Name: "testing",
	Config: *&Config{
		Image:    "bsn/image",
		Tag:      "sha-8475c1f",
		Registry: "docker.pkg.github.com",
	},
}

// Setup resources
func setup() { CreateLocalDB() }

// Remove resources
func shutdown() { RemoveLocalDB() }

func TestMain(m *testing.M) {
	setup()
	defer data.DB.Close()

	code := m.Run()

	shutdown()

	os.Exit(code)
}

func TestTemplateShouldFailIfMissingImageAndName(t *testing.T) {
	newTemplate := ParseTemplate([]byte(templateWithNoImageAndName))
	if newTemplate != nil {
		t.Error("Expected test to fail since template contains no image and no name")
	}
}

func TestTemplateShouldFailIfMissingImage(t *testing.T) {
	newTemplate := ParseTemplate([]byte(templateWithNoImage))
	if newTemplate != nil {
		t.Error("Expected test to fail since template contains no image")
	}
}

func TestTemplateShouldFailIfMissingName(t *testing.T) {
	newTemplate := ParseTemplate([]byte(templateWithNoImage))
	if newTemplate != nil {
		t.Error("Expected test to fail since template contains no name")
	}
}

func TestTemplateShouldNotFailIfContainsNameAndImage(t *testing.T) {
	newTemplate := ParseTemplate([]byte(templateWithNameAndImage))
	if newTemplate == nil {
		t.Error("Expected test to pass since template contains name and image")
	}
}

func TestTemplateDefaultsAreApplied(t *testing.T) {
	testTemplate := &Template{
		Name: "testing",
		Config: *&Config{
			Image: "bsn/image",
		},
	}

	SetTemplateDefaults(testTemplate)

	if testTemplate.Config.Tag == "" {
		t.Error("Expected tag to be `latest`")
	}

	if testTemplate.Config.Registry == "" {
		t.Error("Expected registry to be set to be `docker.io`")
	}
}

func TestTemplateDefaultsAreNotAppliedIfAlreadySet(t *testing.T) {
	SetTemplateDefaults(testTemplate)

	if testTemplate.Config.Tag != "sha-8475c1f" {
		t.Error("Expected tag to be `sha-8475c1f`")
	}

	if testTemplate.Config.Registry != "docker.pkg.github.com" {
		t.Error("Expected registry to be set to be `docker.pkg.github.com`")
	}
}

func TestSaveTemplateToDatastore(t *testing.T) {
	err := SaveDeployment(testTemplate)
	if err != nil {
		t.Errorf("Got error when saving template %s", err.Error())
	}
}

func TestFindTemplate(t *testing.T) {
	tmpl := FindDeployment(testTemplate.Name)

	if tmpl == nil {
		t.Error("Expected template got `nil`")
	}

	if tmpl.Name != testTemplate.Name {
		t.Errorf("Expected template name to be `%s` but got `%s`", testTemplate.Name, tmpl.Name)
	}

	if tmpl.Config.Image != testTemplate.Config.Image {
		t.Errorf("Expected template image to be `%s` but got `%s`", testTemplate.Config.Image, tmpl.Config.Image)
	}

	if tmpl.Config.Tag != testTemplate.Config.Tag {
		t.Errorf("Expected template tag to be `%s` but got `%s`", testTemplate.Config.Tag, tmpl.Config.Tag)
	}

	if tmpl.Config.Registry != testTemplate.Config.Registry {
		t.Errorf("Expected template registry to be `%s` but got `%s`", testTemplate.Config.Registry, tmpl.Config.Registry)
	}
}

// Create & Setup local db
func CreateLocalDB() {
	_, err := data.NewDB("krane.db")
	if err != nil {
		panic("Unable start db")

	}
	err = data.SetupDB()
	if err != nil {
		panic("Unable to start db")
	}
}

// Remove db
func RemoveLocalDB() { os.RemoveAll("./db") }