# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: unary.proto
# Protobuf Python Version: 4.25.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0bunary.proto\x12\x0fshoppingService\"\x14\n\x04Item\x12\x0c\n\x04name\x18\x01 \x01(\t\"<\n\x0e\x43lientItemList\x12*\n\x0b\x63lientItems\x18\x02 \x03(\x0b\x32\x15.shoppingService.Item\"O\n\x0eServerItemList\x12*\n\x0bserverItems\x18\x03 \x03(\x0b\x32\x15.shoppingService.Item\x12\x11\n\ttimestamp\x18\x04 \x01(\t2W\n\x05Unary\x12N\n\x08GetOrder\x12\x1f.shoppingService.ClientItemList\x1a\x1f.shoppingService.ServerItemList\"\x00\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'unary_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  DESCRIPTOR._options = None
  _globals['_ITEM']._serialized_start=32
  _globals['_ITEM']._serialized_end=52
  _globals['_CLIENTITEMLIST']._serialized_start=54
  _globals['_CLIENTITEMLIST']._serialized_end=114
  _globals['_SERVERITEMLIST']._serialized_start=116
  _globals['_SERVERITEMLIST']._serialized_end=195
  _globals['_UNARY']._serialized_start=197
  _globals['_UNARY']._serialized_end=284
# @@protoc_insertion_point(module_scope)