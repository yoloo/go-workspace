#! /usr/bin/python
# -*- coding: gb2312 -*-
################################
# xml::XMLParser generator
#
# usage:
# xmlpg.py --help
###############################

import sys
import os
import getopt
import xml.etree.ElementTree as ET

######################################################
#
# Helper Functions
#
######################################################

def printline(file, depth, text):
	for i in range(0, depth):
		print >> file, '\t',
		print  >> sys.stdout, '\t',
	print >> file, text
	print  >> sys.stdout, text

def NodeAsContainer(node):
	container = ''
	key = ''
	keytype = ''
	for attr in node.attrib:
		if attr in ('container_', 'cont_'):
			container = node.attrib[attr]
		elif attr in ('key_', 'key__'):
			key = node.attrib[attr]

	if key != '':
		for attr_ in node.attrib:
			if attr_ == key:
				keytype = node.attrib[attr_]
				break

	return container, key, keytype

def CreateNameFunc(name):
	pos = name.find('_')
	while pos != -1:
		if len(name) > pos + 1:
			name = name[0:pos] + name[pos + 1: pos + 2].upper() + name[pos + 2:]
		else:
			name = name[0:pos]
		pos = name.find('_')
	return name

def Node2Struct(node):
	return node.tag[0:1].upper() + CreateNameFunc(node.tag[1:])

def AttrVarName(name):
	return name[0:1].upper() + CreateNameFunc(name[1:])

def NodeHasContent(node):
	return node.text and node.text.strip()

def GoType(t):
	if (t == 'DWORD'):
		return 'uint32'
	elif (t == 'QWORD'):
		return 'uint64'
	elif (t == 'double'):
	    return 'float64'
	elif (t in ('bool', 'int', 'uint', 'float')):
		return t
	return 'string'


def IsKeyword(word):
	return word in ('container_', 'cont_', 'key_', 'var_', 'file_', 'index_')

def IsSeqCont(cont):
	return cont in ('vector', 'list', 'dequeue')

def IsMapCont(cont):
	return cont in ('map', 'multimap')


######################################################
#
# Generate the go file with config struct
#
######################################################
class GoFileGenerator(object):

	def __init__(self, arg):
		super(GoFileGenerator, self).__init__()
		self.xmlpg = arg

	def gen(self, tree, file=sys.stdout):
		self.tree = tree
		if self.xmlpg.package:
			printline(file, 0, 'package %s' % self.xmlpg.package)
			self.genHeader(tree.getroot(), 0, file)
		else:
			printline(file, 0, 'package main')
			self.genHeader(tree.getroot(), 0, file)

	def genHeader(self, node, depth, file=sys.stdout):
		container, key, keytype = NodeAsContainer(node)

		if IsMapCont(container) or IsSeqCont(container) :
			if depth == 0:
				printline(file, depth, 'type %s []struct{' % Node2Struct(node))
			else:
				printline(file, depth, '%s []struct{' % Node2Struct(node))
		elif depth == 0:
			printline(file, depth, 'type %s struct{' % Node2Struct(node))
		else:
			printline(file, depth, '%s struct{' % Node2Struct(node))
		nodetag = node.tag

		#		if len(node.attrib) > 0:
#			printline(file, depth + 1, '')

		#members
		for attr in node.attrib:
			if not IsKeyword(attr):
				printline(file, depth + 1, '%s %s `xml:"%s,attr"`' % (AttrVarName(attr), GoType(node.attrib[attr]), attr))

		for child in node:
			self.genHeader(child, depth + 1, file)
		if depth == 0 :
			printline(file, depth, '}')
		else:
			printline(file, depth, '} `xml:"%s"`' % (nodetag))

#			printline(file, depth, '')

######################################################
#
# Generator
#
######################################################
class XMLParserGenerator(object):
	"""docstring for XMLParserGenerator"""
	def __init__(self):
		super(XMLParserGenerator, self).__init__()
		self.header = 0
		self.cpp = 0
		self.all = 0
		self.inline = 0
		self.outputfile = ""
		self.package = ""
		self.schemas = []

	def schemaFileName(self, schema):
		return os.path.splitext(os.path.basename(schema))[0]

	def outputFileName(self, schema, ext):
		return os.path.join("path...", self.schemaFileName(schema) + '.' + ext)

	def run(self):
		for i in range(0, len(self.schemas)):
			self.generate(self.schemas[i])

	def generate(self, filename):
		try:
			tree = ET.ElementTree(file = filename)

			file = sys.stdout

			self.gg = GoFileGenerator(self)

			if (self.outputfile):
				file = open(self.outputfile + '.go', 'w+')
			self.gg.gen(tree, file)

		except IOError, arg:
			print
			arg




######################################################
#
# main
#
######################################################
xmlpg = XMLParserGenerator()

def Usage():
	print("usage:", sys.argv[0], "[-i -o -file] schemal.xml schema2.xml ...")
	print("options: ")
	print("	-h --help				help")
	print("	-o file --output=file	output to the file")
	print("	-p --package			package")

if len(sys.argv) == 1:
	Usage()
	sys.exit()

try:
	opts, args = getopt.getopt(sys.argv[1:], 'ho:p:', ['help', 'output=','package'])
except getopt.GetoptError:
	Usage()
	sys.exit()

for o, a in opts:
	if o in ("-h", "--help"):
		Usage()
		sys.exit()
	elif o in ("-o", "--output"):
		xmlpg.outputfile = a
	elif o in ("-p", "--package"):
		xmlpg.package = a

xmlpg.schemas = args
if len(xmlpg.schemas) > 0:
	xmlpg.run()
else:
	Usage()
