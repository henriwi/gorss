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

gorss.controller("HomeCtrl", ["$scope", "Feed", "$http", "$window", function($scope, Feed, $http, $window) {
  $scope.loading = true;
	$scope.feeds = Feed.query(function() {
    $scope.loading = false;
  });

	$scope.update = function() {
    $scope.loading = true;
    $scope.feeds = Feed.query(function(data){
      $scope.loading = false;
    });
	};

	$scope.addFeed = function() {
    $http({
      method: "POST",
      url: "/api/feed",
      data: {url: $scope.url}
    }).success(function(data) {
      $scope.feeds.push(data);
    }).error(function(data) {
      $(".alert-danger").fadeIn(500);
      setTimeout(function() {
        $scope.$apply(function(){
          // $scope.error = false;
          $(".alert-danger").fadeOut(500);
        });
      }, 2000);
    });

    // Feed.save("url:" + $scope.url);
    // $scope.feeds.push({Title: $scope.url})

    $scope.url = "";
	}

  $scope.markUnread = function(link, feedIndex, itemIndex) {
    $scope.feeds[feedIndex].Items[itemIndex].Read = true

    $http({
      method: "POST",
      url: "/api/feed/"+feedIndex+"/item/"+itemIndex+"/read",
      data: {url: $scope.url}
    });

    $window.open(link);
  };

}]);

gorss.controller("FeedCtrl", ["$scope", "Feed", function($scope, Feed, $routeParams) {
	var url = $routeParams.url;
	$scope.feeds = Feed.query();

  $scope.loading = true;

	$scope.update = function() {
		$scope.feeds = Feed.query();
    $scope.loading = false;
	};

}]);