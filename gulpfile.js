var gulp = require('gulp');
var favicon = require ('gulp-real-favicon');
var fs = require('fs');
// Using until gulp v4 is released
var runSequence = require('run-sequence');

var faviconData = 'asset/dynamic/favicon/data.json';

// SASS Task
gulp.task('sass', function() {
    var sass = require('gulp-sass');
    return gulp.src('asset/dynamic/sass/**/*.scss')
        // Available for outputStyle: expanded, nested, compact, compressed
        .pipe(sass({outputStyle: 'expanded'}).on('error', sass.logError))
        .pipe(gulp.dest('asset/static/css/'));
});

// JavaScript Task
gulp.task('javascript', function() {
	var concat = require('gulp-concat');
	return gulp.src('asset/dynamic/js/*.js')
		.pipe(concat('all.js'))
		.pipe(gulp.dest('asset/static/js/'));
});

// jQuery Task
gulp.task('jquery', function() {
	return gulp.src('bower_components/jquery/dist/jquery.min.*')
		.pipe(gulp.dest('asset/static/js/'));
});

// Bootstrap Task
gulp.task('bootstrap', function() {
	gulp.src('bower_components/bootstrap/dist/css/bootstrap-theme.min.*')
		.pipe(gulp.dest('asset/static/css/'));
	gulp.src('bower_components/bootstrap/dist/css/bootstrap.min.*')
		.pipe(gulp.dest('asset/static/css/'));
	gulp.src('bower_components/bootstrap/dist/fonts/*')
		.pipe(gulp.dest('asset/static/fonts/'));
	return gulp.src('bower_components/bootstrap/dist/js/bootstrap.min.js')
		.pipe(gulp.dest('asset/static/js/'));
});

// Underscore Task
gulp.task('underscore', function() {
	return gulp.src('bower_components/underscore/underscore-min.*')
		.pipe(gulp.dest('asset/static/js/'));
});

// Favicon Generation and Injection Task
gulp.task('favicon', function() {
	runSequence('favicon-generate', 'favicon-inject');
});

// Generate the icons. This task takes a few seconds to complete.
// You should run it at least once to create the icons. Then,
// you should run it whenever RealFaviconGenerator updates its
// package (see the favicon-update task below).
gulp.task('favicon-generate', function(done) {
	var favColor = '#525252';
	favicon.generateFavicon({
		masterPicture: 'asset/dynamic/favicon/logo.png',
		dest: 'asset/static/favicon/',
		iconsPath: '/static/favicon/',
		design: {
			ios: {
				pictureAspect: 'backgroundAndMargin',
				backgroundColor: favColor,
				margin: '14%'
			},
			desktopBrowser: {},
			windows: {
				pictureAspect: 'noChange',
				backgroundColor: favColor,
				onConflict: 'override'
			},
			androidChrome: {
				pictureAspect: 'noChange',
				themeColor: favColor,
				manifest: {
					name: 'Blueprint',
					display: 'browser',
					orientation: 'notSet',
					onConflict: 'override',
					declared: true
				}
			},
			safariPinnedTab: {
				pictureAspect: 'silhouette',
				themeColor: favColor
			}
		},
		settings: {
			scalingAlgorithm: 'Mitchell',
			errorOnImageTooSmall: false
		},
		versioning: {
			paramName: 'v1.0',
			paramValue: '3eepn6WlLO'
		},
		markupFile: faviconData
	}, function() {
		done();
	});
});

// Inject the favicon markups in your HTML pages. You should run
// this task whenever you modify a page. You can keep this task
// as is or refactor your existing HTML pipeline.
gulp.task('favicon-inject', function() {
	return gulp.src(['view/partial/favicon.tmpl'])
		.pipe(favicon.injectFaviconMarkups(JSON.parse(fs.readFileSync(faviconData)).favicon.html_code))
		.pipe(gulp.dest('view/partial/'));
});

// Check for updates on RealFaviconGenerator (think: Apple has just
// released a new Touch icon along with the latest version of iOS).
// Run this task from time to time. Ideally, make it part of your
// continuous integration system.
gulp.task('favicon-update', function(done) {
	var currentVersion = JSON.parse(fs.readFileSync(faviconData)).version;
	return favicon.checkForUpdates(currentVersion, function(err) {
		if (err) {
			throw err;
		}
	});
});

// Watch
gulp.task('watch', function() {
    gulp.watch('asset/dynamic/sass/**/*.scss', ['sass']);
	gulp.watch('asset/dynamic/js/*.js', ['javascript']);
});

// Init - every task
gulp.task('init', ['sass', 'javascript', 'jquery', 'bootstrap', 'underscore', 'favicon']);

// Default - only run the tasks that change often
gulp.task('default', ['sass', 'javascript']);