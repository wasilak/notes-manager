/* jslint node: true */
"use strict";

function NoteRenderedCtrl($rootScope, $scope, $stateParams, ApiService) {
  var vm = this;

  vm.uuid = $stateParams.uuid;

  vm.note = {};

  vm.inputText = '';

  ApiService.getNote(vm.uuid).then(function(result) {
    vm.note = result;
    vm.outputText = marked(result.content);
  });
}

NoteRenderedCtrl.resolve = {
};

angular.module("app").controller("NoteRenderedCtrl", NoteRenderedCtrl);
