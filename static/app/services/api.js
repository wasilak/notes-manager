function ApiService($http, API, $rootScope, APP_SETTINGS) {
  var getNote = function(uuid) {

    var url = API.urls.note.replace("{{uuid}}", uuid);

    return $http({
        cache: false,
        url: url,
        method: 'GET',
        params: {}
      })
      .then(function(response) {
        return {
          response: response.data.data,
          success: true
        };
      },
      function(response) {
        console.error('Getting note failed!');
        return {
          success: false,
          error: response.data.error,
          status: response.status,
          statusText: response.statusText
        };
      }
    );
  };

  var getList = function(filter, sort, tags) {

    if (!filter) {
      filter = "";
    }

    if (!sort) {
      sort = "";
    }

    if (!tags) {
      tags = [];
    }

    var url = API.urls.list.replace("{{filter}}", filter);

    return $http({
        cache: false,
        url: url,
        method: 'GET',
        params: {
          sort: sort,
          tags: tags.join(",")
        }
      })
      .then(function(response) {
        return {
          response: response.data.data,
          success: true
        };
      },
      function(response) {
        console.error('Getting notes list failed!');
        console.log(response.data);
        return {
          success: false,
          error: response.data.error,
          status: response.status,
          statusText: response.statusText
        };
      }
    );
  };

  var saveNote = function(note) {

    var url = API.urls.note.replace("{{uuid}}", note.response.id);

    return $http({
        cache: false,
        url: url,
        method: 'POST',
        data: {
          note: note.response
        }
      })
      .then(function(response) {
        return {
          response: response.data.data,
          success: true
        };
      },
      function(response) {
        console.error('Saving note failed!');
        return {
          success: false,
          error: response.data.error,
          status: response.status,
          statusText: response.statusText
        };
      }
    );
  };

  var createNote = function(note) {

    var url = API.urls.new;

    return $http({
        cache: false,
        url: url,
        method: 'POST',
        data: {
          note: note.response
        }
      })
      .then(function(response) {
        return {
          response: response.data.data,
          success: true
        };
      },
      function(response) {
        console.error('Creating note failed!');
        return {
          success: false,
          error: response.data.error,
          status: response.status,
          statusText: response.statusText
        };
      }
    );
  };

  var deleteNote = function(uuid) {

    var url = API.urls.delete.replace("{{uuid}}", uuid);

    return $http({
        cache: false,
        url: url,
        method: 'DELETE',
        data: {}
      })
      .then(function(response) {
        return {
          response: response.data.data,
          success: true
        };
      },
      function(response) {
        console.error('Creating note failed!');
        return {
          success: false,
          error: response.data.error,
          status: response.status,
          statusText: response.statusText
        };
      }
    );
  };

  var getTags = function(query) {

    var url = API.urls.tags.replace("{{query}}", query);

    return $http({
        cache: false,
        url: url,
        method: 'GET',
        params: {}
      });
  };

  return {
    getNote: getNote,
    saveNote: saveNote,
    createNote: createNote,
    deleteNote: deleteNote,
    getList: getList,
    getTags: getTags,
  };
}

angular.module("app").factory("ApiService", ApiService);
