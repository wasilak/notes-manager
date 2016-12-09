/**
* node.js API server
**/

/* jshint node: true */
/* jshint esversion: 6 */

"use strict";

// including Express.js module
const express = require('express');
const morgan = require('morgan');
const path = require('path');
const lowdb = require('lowdb');
const uuidV1 = require('uuid/v1');

const bodyParser = require('body-parser');

let app = express();

app.use(bodyParser.json());       // to support JSON-encoded bodies
app.use(bodyParser.urlencoded({     // to support URL-encoded bodies
  extended: true
}));

const db = lowdb('db.json');

db.defaults({ notes: [] }).value();

let noteObject = {
  id: null,
  content: '',
  title: '',
  created: '',
  updated: ''
};

let appPort = 0;
if (process.env.QPP_PORT) {
    appPort = process.env.QPP_PORT;
} else if(process.env.PORT) {
    appPort = process.env.PORT;
} else {
    appPort = 5000;
}

app.set('port', appPort);

app.use(morgan('dev'));

// path to static assets (CSS etc.)
app.use(express.static('public'));

app.get('/api/list', (req, res) => {

  let notes = db.get('notes');

  // cloning object
  // let newNote = Object.assign({}, noteObject);
  // newNote.content = 'dkjsahdasgfuysdgfuyds';
  // newNote.id = uuidV1();

  // let date = {
  //   created: Math.floor(Date.now() / 1000),
  //   updated: Math.floor(Date.now() / 1000)
  // };

  // newNote.created = date.created;
  // newNote.updated = date.updated;

  // console.log(noteObject, newNote);

  // notes.push(newNote).value();

  res.setHeader('Content-Type', 'application/json');
  res.send(JSON.stringify(notes));
});

app.get('/api/note/:uuid', (req, res) => {

    let uuid = req.params.uuid;
    let note = db.get('notes').find({ id: uuid }).value();

    res.setHeader('Content-Type', 'application/json');
    res.send(JSON.stringify(note));
});

app.post('/api/note/:uuid', (req, res) => {
  let uuid = req.params.uuid;
  let note = req.body.note;
  let notes = db.get('notes');

  note = notes.find({ id: uuid }).assign({
    content: note.content,
    title: note.title,
    updated: Math.floor(Date.now() / 1000)
  }).value();

  res.setHeader('Content-Type', 'application/json');
  res.send(JSON.stringify(note));
});

// catching all routes with single page AngularJS app.
// AngularJS will take care of the routing.
app.get('*', (req, res) => {
     res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// server init on custom port
let server = app.listen(app.get('port'), () => {
    let host = server.address().address;
    let port = server.address().port;

});
