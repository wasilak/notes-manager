function config($httpProvider, $compileProvider, $stateProvider, $urlRouterProvider) {
  // ui.router
  $stateProvider
    // base parent state
    .state('parent', {
      abstract: true,
      template: '<ui-view/>',
      views: {
        'menu': {
          templateUrl: 'app/views/menu.html',
          controller: 'MenuCtrl as vm'
        }
      }
    })
    .state('parent.list', {
      url: '/',
      data: {
        title: 'List',
      },
      views: {
        '@': {
          templateUrl: 'app/views/list.html',
          controller: 'ListCtrl as vm'
        }
      },
      resolve: ListCtrl.resolve
    })
    .state('parent.new', {
      url: '/note/new',
      data: {
        title: 'New note',
      },
      views: {
        '@': {
          templateUrl: 'app/views/note.html',
          controller: 'NewCtrl as vm'
        }
      },
      resolve: NewCtrl.resolve
    })
    .state('parent.note', {
      url: '/note/:uuid',
      data: {
        title: 'Note',
      },
      views: {
        '@': {
          templateUrl: 'app/views/note.html',
          controller: 'NoteCtrl as vm'
        }
      },
      resolve: NoteCtrl.resolve
    })
    .state('parent.list.note', {
      url: 'list/:uuid',
      data: {
        title: 'List :: Note',
      },
      views: {
        'content@parent.list': {
          templateUrl: 'app/views/noteRendered.html',
          controller: 'NoteRenderedCtrl as vm'
        }
      },
      resolve: NoteRenderedCtrl.resolve
    })
    ;

  // /emails -> /inbox
  // Automated redirects
  $urlRouterProvider.otherwise('/');

  // batches $digest cycles for $http calls
  // that resolve within 10ms of eachother
  $httpProvider.useApplyAsync(true);

  // BEFORE: <div ng-controller="ListCtrl as list" class="ng-controller ng-binding"></div>
  //         angular.element('.myClass').scope();
  //
  // AFTER:  <div ng-controller="ListCtrl as list"></div>
  //         - Doesn't add unnecessary class names
  //         - Doesn't bind .scope() / .getIsolateScope() data to each element
  $compileProvider.debugInfoEnabled(true);

}

function run($rootScope, APP_SETTINGS) {
  var page = {
    appName: APP_SETTINGS.name,
    setTitle: function (title) {
      this.title = APP_SETTINGS.name + ' :: ' + title;
    }
  };

  function setTitle(event, state) {
    page.setTitle(state && state.data ? state.data.title : '');
  }

  // exports
  $rootScope.page = page;
  $rootScope.$on('$stateChangeSuccess', setTitle);
}

angular
  .module('app')
  .run(run)
  .config(config);
