var gorss = angular.module('gorss', ["feedService"]);

gorss.config(function($locationProvider, $routeProvider) {
  // $locationProvider.html5Mode(true);
  $routeProvider.
    when('/feed/:url', {
      templateUrl: 'templates/feed.html',
      controller: 'FeedCtrl'
    }).
    when('/', {
      templateUrl: 'templates/home.html',
      controller: "HomeCtrl"
    }).
    otherwise({redirectTo: '/404'});
});

var feedService = angular.module("feedService", ["ngResource"]);

feedService.factory('Feed', ['$resource',
  function($resource){
    return $resource('api/feed', {}, {
      query: {method:'GET', isArray:true}
    });
  }]);

gorss.controller("HomeCtrl", ["$scope", "Feed", "$http", function($scope, Feed, $http) {
	$scope.feeds = Feed.query();

	$scope.update = function() {
		$scope.feeds = Feed.query();
	};

	$scope.addFeed = function() {
    $http({
      method: "POST",
      url: "/api/feed",
      data: {url: $scope.url}
    }).success(function() {
      $scope.feeds = Feed.query();
    });

    // Feed.save("url:" + $scope.url);
    // $scope.feeds.push({Title: $scope.url})

    $scope.url = "";
	}

}]);

gorss.controller("FeedCtrl", ["$scope", "Feed", function($scope, Feed, $routeParams) {
	var url = $routeParams.url;
	$scope.feeds = Feed.query();

	$scope.update = function() {
		$scope.feeds = Feed.query();
	};

}]);