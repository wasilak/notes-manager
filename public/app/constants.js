angular.module("app").constant("APP_SETTINGS", {
  name: "Notes Manager"
})
.constant("API", {
  urls: {
    note: "/api/note/{{uuid}}",
    new: "/api/note/new",
    delete: "/api/note/delete",
    list: "/api/list"
  }
});
