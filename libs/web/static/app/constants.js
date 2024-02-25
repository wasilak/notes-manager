/* global angular */

angular.module("app").constant("APP_SETTINGS", {
  name: "Notes Manager"
})
  .constant("API", {
    urls: {
      note: "/api/note/{{uuid}}",
      new: "/api/note/",
      list: "/api/list/",
      tags: "/api/tags/",
      aiReWriteNote: "/api/ai/rewrite/",
      aiEnabled: "/api/ai/enabled/"
    }
  });
