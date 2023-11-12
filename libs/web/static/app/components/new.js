/* jslint node: true */
/* global angular */
/* global marked */

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

      vm.createNote = function () {
        vm.loader = true;
        ApiService.createNote(vm.note).then(function (result) {
          $rootScope.notifications.push('Note created');
          vm.loader = false;
          $state.go('list_note', { uuid: result.response.id }, { reload: true });
        });
      };

      vm.createButtonDisabled = function () {
        return vm.note.response.content.length === 0 || vm.note.response.title.length === 0;
      };

      vm.aiRewriteButtonDisabled = function () {
        return vm.note.response.content.length === 0;
      };

      vm.loadItems = function (query) {
        return ApiService.getTags(query);
      };

      vm.aiReWriteNote = function () {
        vm.loader = true;
        ApiService.aiReWriteNote(vm.note).then(function (result) {
          console.log(result)
          if (result.response.rewritten.error) {
            $rootScope.notifications.push('AIRewrite error: ' + result.response.rewritten.error);
          } else {
            vm.note.response.title = result.response.rewritten.title;
            vm.note.response.content = result.response.rewritten.content;
            vm.note.response.tags = result.response.rewritten.tags;
          }
          vm.loader = false;
        });
      };

      // eslint-disable-next-line no-unused-vars
      $scope.$watch('$ctrl.note', function (current, original) {
        vm.errorMessage = false;
        try {
          vm.outputText = marked.parse(current.response.content);
        } catch (err) {
          vm.errorMessage = err.message;
        }

        $rootScope.$broadcast('currentNote', current);
      }, true);
    },
    templateUrl: "/static/app/views/note.html"
  }
);
