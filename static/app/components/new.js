/* jslint node: true */
"use strict";

angular.module("app").component("new", 
  {
    controller: function ($scope, $rootScope, $stateParams, ApiService, $state) {
      var vm = this;

      vm.note = {
        success: true,
        response: {
          id: null,
          content: '',
          title: '',
          created: '',
          updated: ''
        }
      };

      vm.createNote = function() {
        ApiService.createNote(vm.note).then(function(result) {
          $rootScope.notifications.push('Note created');
          $state.go('list_note', {uuid: result.response.id}, {reload: true});
        });
      };

      $scope.$watch('$ctrl.note', function(current, original) {
        vm.errorMessage = false;
        try {
          vm.outputText = marked(current.response.content);
        } catch (err) {
          vm.errorMessage = err.message;
        }

        $rootScope.$broadcast('currentNote', current);
      }, true);
    },
    templateUrl: "/static/app/views/note.html"
  }
);
