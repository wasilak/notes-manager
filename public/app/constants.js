angular.module("app").constant("APP_SETTINGS", {
  name: "Notes Manager"
})
.constant("API", {
  urls: {
    note: "/api/note/{{uuid}}",
    list: "/api/list"
  }
});
