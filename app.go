package main

import (
	"context"

	"github.com/madeindra/interview-app/internal/database"
	"github.com/madeindra/interview-app/internal/elevenlabs"
	"github.com/madeindra/interview-app/internal/model"
	"github.com/madeindra/interview-app/internal/openai"
)

// App struct
type App struct {
	ctx    context.Context
	model  *model.Model
	oaiAPI *openai.OpenAI
	elAPI  *elevenlabs.ElevenLab
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db := database.New()
	a.model = model.New(db)

	a.oaiAPI = openai.New()
	a.elAPI = elevenlabs.New()
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}
