/* jslint node: true */
/* global angular */
"use strict";

angular.module("app").component("list",
  {
    controller: function ListCtrl($rootScope, ApiService, $scope) {
      var vm = this;

      vm.list = {
        success: true
      };

      vm.sort = "updated:desc";

      vm.tags = [];

      $rootScope.$on('currentNote', function (event, note) {
        vm.note = note;
      });

      vm.updateList = function () {
        ApiService.getList(vm.listFilter, vm.sort, vm.tags).then(function (result) {
          vm.list = result;
        });
      };

      vm.setSort = function () {
        vm.updateList();
      };

      vm.clearSearch = function () {
        vm.listFilter = "";
        vm.search();
      };

      vm.search = function () {
        vm.updateList();
      };

      vm.loadItems = function (query) {
        return ApiService.getTags(query);
      };

      // eslint-disable-next-line no-unused-vars
      $scope.$watch('$ctrl.tags', function (current, original) {
        vm.search();
      });

      vm.loadItems();

      vm.listFilter = "";

    },
    templateUrl: "/static/app/views/list.html"
  });
