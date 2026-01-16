package gplay

import "context"

func Default() *Client { return DefaultClient }

func FetchApp(ctx context.Context, opts AppOptions) (App, error) { return DefaultClient.App(ctx, opts) }
func FetchList(ctx context.Context, opts ListOptions) ([]App, error) {
	return DefaultClient.List(ctx, opts)
}
func FetchSearch(ctx context.Context, opts SearchOptions) ([]App, error) {
	return DefaultClient.Search(ctx, opts)
}
func FetchDeveloper(ctx context.Context, opts DeveloperOptions) ([]App, error) {
	return DefaultClient.Developer(ctx, opts)
}
func FetchSuggest(ctx context.Context, opts SuggestOptions) ([]string, error) {
	return DefaultClient.Suggest(ctx, opts)
}
func FetchReviews(ctx context.Context, opts ReviewsOptions) (ReviewsResult, error) {
	return DefaultClient.Reviews(ctx, opts)
}
func FetchSimilar(ctx context.Context, opts SimilarOptions) ([]App, error) {
	return DefaultClient.Similar(ctx, opts)
}
func FetchPermissions(ctx context.Context, opts PermissionsOptions) (PermissionsResult, error) {
	return DefaultClient.Permissions(ctx, opts)
}
func FetchDataSafety(ctx context.Context, opts DataSafetyOptions) (DataSafetyResult, error) {
	return DefaultClient.DataSafety(ctx, opts)
}
func FetchCategories(ctx context.Context, opts CategoriesOptions) ([]string, error) {
	return DefaultClient.Categories(ctx, opts)
}
