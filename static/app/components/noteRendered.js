/* jslint node: true */
"use strict";

angular.module("app").component("noteRendered", 
  {
    bindings: {
      note: '<'
    },
    controller: function ($scope, $rootScope, $stateParams, ApiService, $state) {
      var vm = this;

      $scope.$watch('$ctrl.note.response.content', function(current, original) {
        vm.errorMessage = false;
        try {
          vm.outputText = marked(current);
        } catch (err) {
          vm.errorMessage = err.message;
        }
      });

      vm.deleteNote = function() {
        var confirmed = confirm("Are you sure?");

        if (confirmed) {
          ApiService.deleteNote(vm.note.response.id).then(function(result) {
              $rootScope.notifications.push('Note deleted');
              $state.go('list', {}, {reload: true});
          });
        }
      };

    },
    templateUrl: "/static/app/views/noteRendered.html"
  }
);
