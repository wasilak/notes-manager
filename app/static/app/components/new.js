/* jslint node: true */
"use strict";

angular.module("app").component("new", 
  {
    controller: function ($scope, $rootScope, $stateParams, ApiService, $state) {
      var vm = this;

      vm.loader = false;

      vm.note = {
        success: true,
        response: {
          id: null,
          content: '',
          title: '',
          created: null,
          updated: null,
          tags: []
        }
      };

      vm.createNote = function() {
        vm.loader = true;
        ApiService.createNote(vm.note).then(function(result) {
          $rootScope.notifications.push('Note created');
          vm.loader = false;
          $state.go('list_note', {uuid: result.response.id}, {reload: true});
        });
      };

      vm.createButtonDisabled = function() {
        return vm.note.response.content.length == 0 || vm.note.response.title.length == 0;
      };

      vm.loadItems = function(query) {
        return ApiService.getTags(query);
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
