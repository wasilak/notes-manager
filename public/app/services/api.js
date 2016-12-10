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
        return response.data;
      },
      function(response) {
        console.error('Getting note failed!');
        console.error(response);
        return {
          success: false,
          error: response
        };
      }
    );
  };

  var getList = function() {

    var url = API.urls.list;

    return $http({
        cache: false,
        url: url,
        method: 'GET',
        params: {}
      })
      .then(function(response) {
        return response.data;
      },
      function(response) {
        console.error('Getting notes list failed!');
        console.error(response);
        return {
          success: false,
          error: response
        };
      }
    );
  };

  var saveNote = function(note) {

    var url = API.urls.note.replace("{{uuid}}", note.id);

    return $http({
        cache: false,
        url: url,
        method: 'POST',
        data: {
          note: note
        }
      })
      .then(function(response) {
        return response.data;
      },
      function(response) {
        console.error('Saving note failed!');
        console.error(response);
        return {
          success: false,
          error: response
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
          note: note
        }
      })
      .then(function(response) {
        return response.data;
      },
      function(response) {
        console.error('Creating note failed!');
        console.error(response);
        return {
          success: false,
          error: response
        };
      }
    );
  };

  var deleteNote = function(uuid) {

    var url = API.urls.delete;

    return $http({
        cache: false,
        url: url,
        method: 'POST',
        data: {
          uuid: uuid
        }
      })
      .then(function(response) {
        return response.data;
      },
      function(response) {
        console.error('Creating note failed!');
        console.error(response);
        return {
          success: false,
          error: response
        };
      }
    );
  };

  return {
    getNote: getNote,
    saveNote: saveNote,
    createNote: createNote,
    deleteNote: deleteNote,
    getList: getList
  };
}

angular.module("app").factory("ApiService", ApiService);
