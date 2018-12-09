/* jslint node: true */
"use strict";

function ListCtrl($rootScope, ApiService) {
  var vm = this;

  vm.list = [];

  $rootScope.$on('currentNote', function(event, note) {
    vm.note = note;
  });

  vm.search = function() {
    ApiService.getList(vm.listFilter, vm.sort).then(function(result) {
       vm.list = result;
    });
  };

  vm.setSort = function() {
    vm.search();
  };

  vm.search();
}

ListCtrl.resolve = {
  notes: function($stateParams, ApiService, $rootScope) {
    return ApiService.getList();
  }
};

angular.module("app").controller("ListCtrl", ListCtrl);
