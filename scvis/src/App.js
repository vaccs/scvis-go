import React from 'react';
import ControlPanel from './ControlPanel';
import Options from './Options';
import About from './About';
import {
  FileCloseEvent, FileOpenEvent, FileReadEvent, FileWriteEvent, FreeEvent, MallocEvent, FunctionInvocationEvent, LineChangeEvent,
  OutputEvent, RegisterEvent, ReturnEvent, ReturnAddressEvent, VariableAccessEvent
} from './Event';
import { parseSectionInfo } from './Sections';
import { ASMSourceMap } from './ASMSource';
import { useEffect, useState } from 'react';
import axios from 'axios';

function makeElementObj(element) {
  var elementObj;

  if (element.VACCS_FILE_CLOSE) {
    elementObj = new FileCloseEvent(element.VACCS_FILE_CLOSE.Event_num,
      element.VACCS_FILE_CLOSE.Line_num,
      element.VACCS_FILE_CLOSE.Path,
      element.VACCS_FILE_CLOSE.File_descriptor);
  } else if (element.VACCS_FILE_OPEN) {
    elementObj = new FileOpenEvent(element.VACCS_FILE_OPEN.Event_num,
      element.VACCS_FILE_OPEN.Line_num,
      element.VACCS_FILE_OPEN.Path,
      element.VACCS_FILE_OPEN.File_descriptor,
      element.VACCS_FILE_OPEN.Flags,
      element.VACCS_FILE_OPEN.Mode,
      element.VACCS_FILE_OPEN.Is_sym_link);
  } else if (element.VACCS_FILE_READ) {
    elementObj = new FileReadEvent(element.VACCS_FILE_READ.Event_num,
      element.VACCS_FILE_READ.Line_num,
      element.VACCS_FILE_READ.Path,
      element.VACCS_FILE_READ.File_descriptor,
      element.VACCS_FILE_READ.Valid);
  } else if (element.VACCS_FILE_WRITE) {
    elementObj = new FileWriteEvent(element.VACCS_FILE_WRITE.Event_num,
      element.VACCS_FILE_WRITE.Line_num,
      element.VACCS_FILE_WRITE.Path,
      element.VACCS_FILE_WRITE.File_descriptor,
      element.VACCS_FILE_WRITE.Valid);
  } else if (element.VACCS_FREE) {
    elementObj = new FreeEvent(element.VACCS_FREE.Event_num,
      element.VACCS_FREE.Line_num,
      element.VACCS_FREE.Path,
      element.VACCS_FREE.Address);
  } else if (element.VACCS_MALLOC) {
    elementObj = new MallocEvent(element.VACCS_MALLOC.Event_num,
      element.VACCS_MALLOC.Line_num,
      element.VACCS_MALLOC.Path,
      element.VACCS_MALLOC.Address,
      element.VACCS_MALLOC.Size);
  } else if (element.VACCS_FUNCTION_INV) {
    elementObj = new FunctionInvocationEvent(element.VACCS_FUNCTION_INV.Event_num,
      element.VACCS_FUNCTION_INV.Line_num,
      element.VACCS_FUNCTION_INV.Path,
      element.VACCS_FUNCTION_INV.Inv_line_num,
      element.VACCS_FUNCTION_INV.C_inv_path,
      element.VACCS_FUNCTION_INV.Func_name,
      element.VACCS_FUNCTION_INV.Address);
  } else if (element.VACCS_LINE) {
    elementObj = new LineChangeEvent(element.VACCS_LINE.Event_num,
      element.VACCS_LINE.Line_num,
      element.VACCS_LINE.Path);
  } else if (element.VACCS_OUTPUT) {
    elementObj = new OutputEvent(element.VACCS_OUTPUT.Event_num,
      element.VACCS_OUTPUT.Line_num,
      element.VACCS_OUTPUT.Path,
      element.VACCS_OUTPUT.Output);
  } else if (element.VACCS_REGISTER) {
    elementObj = new RegisterEvent(element.VACCS_REGISTER.Event_num,
      element.VACCS_REGISTER.Line_num,
      element.VACCS_REGISTER.Path,
      element.VACCS_REGISTER.Register,
      element.VACCS_REGISTER.Value);
  } else if (element.VACCS_RETURN) {
    elementObj = new ReturnEvent(element.VACCS_RETURN.Event_num,
      element.VACCS_RETURN.Line_num,
      element.VACCS_RETURN.Path);
  } else if (element.VACCS_RETURN_ADDR) {
    elementObj = new ReturnAddressEvent(element.VACCS_RETURN_ADDR.Event_num,
      element.VACCS_RETURN_ADDR.Line_num,
      element.VACCS_RETURN_ADDR.Path,
      element.VACCS_RETURN_ADDR.Dynamic_link,
      element.VACCS_RETURN_ADDR.Return_address,
      element.VACCS_RETURN_ADDR.Function_name);
  } else if (element.VACCS_VAR_ACCESS) {
    elementObj = new VariableAccessEvent(element.VACCS_VAR_ACCESS.Event_num,
      element.VACCS_VAR_ACCESS.Line_num,
      element.VACCS_VAR_ACCESS.Path,
      element.VACCS_VAR_ACCESS.Scope,
      element.VACCS_VAR_ACCESS.Address,
      element.VACCS_VAR_ACCESS.Var_type,
      element.VACCS_VAR_ACCESS.Name,
      element.VACCS_VAR_ACCESS.Points_to,
      element.VACCS_VAR_ACCESS.Value);
  }
  else {
    console.log("Unknown element type: " + element);
    elementObj = null;
  }

  return elementObj;
}

function addToArray(array, element) {
  var elementObj = makeElementObj(element);
  if (array === undefined) {
    return [elementObj];
  } else {
    array.push(elementObj);
    return array;
  }
}

// const jsonData = [{"VACCS_ARCH":{"Event_num":-1,"Arch_type":1,"Heap_start":94400437915648,"Heap_end":140728349331448,"Stack_start":140728349470720,"Stack_end":140728349331456}},{"VACCS_CMD_LINE":{"Event_num":-1,"Cmd_line":"./t1 "}},{"VACCS_BINARY":{"Event_num":-1,"Path":"./t1"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":1,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"#include \u003cstdio.h\u003e"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":2,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"int a;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":3,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"int c = 3;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"int main() {"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":5,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":""}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":6,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   int b;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":7,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":""}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":8,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   int *p = NULL;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":9,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":""}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":10,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   p = \u0026b; "}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":11,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   a = 1;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":12,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   b = 2;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   printf(\"a = %d, b = %d, *p = %d, c = %d\\n\",a, b, *p, c);"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"}"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|Type|Address|Size|Section Name"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8318|28|.interp"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8338|48|.note.gnu.property"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8368|36|.note.gnu.build-id"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb838c|32|.note.ABI-tag"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb83b0|36|.gnu.hash"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|16|0x55db4eeb83d8|168|.dynsym"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|20|0x55db4eeb8480|143|.dynstr"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8510|14|.gnu.version"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8520|48|.gnu.version_r"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|2|0x55db4eeb8550|192|.rela.dyn"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|2|0x55db4eeb8610|24|.rela.plt"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9000|27|.init"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9020|32|.plt"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9040|16|.plt.got"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9050|16|.plt.sec"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9060|331|.text"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb91ac|13|.fini"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeba000|41|.rodata"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeba02c|52|.eh_frame_hdr"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeba060|172|.eh_frame"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eebbdb8|8|.init_array"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eebbdc0|8|.fini_array"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|6|0x55db4eebbdc8|496|.dynamic"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|8|0x55db4eebbfb8|72|.got"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eebc000|20|.data"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|18|0x55db4eebc014|12|.bss"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|43|.comment"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|48|.debug_aranges"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|244|.debug_info"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|160|.debug_abbrev"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|98|.debug_line"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|220|.debug_str"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|84|.debug_line_str"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|15|0x0|912|.symtab"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|19|0x0|477|.strtab"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|21|0x0|362|.shstrtab"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37193,"C_line_num":4,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"nop edx, edi"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37197,"C_line_num":4,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"push rbp"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37198,"C_line_num":4,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov rbp, rsp"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37201,"C_line_num":4,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"sub rsp, 0x10"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37205,"C_line_num":8,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov qword ptr [rbp-0x8], 0x0"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37213,"C_line_num":10,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"lea rax, ptr [rbp-0xc]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37217,"C_line_num":10,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov qword ptr [rbp-0x8], rax"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37221,"C_line_num":11,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov dword ptr [rip+0x2ea9], 0x1"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37231,"C_line_num":12,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov dword ptr [rbp-0xc], 0x2"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37238,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov esi, dword ptr [rip+0x2e94]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37244,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov rax, qword ptr [rbp-0x8]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37248,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov ecx, dword ptr [rax]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37250,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov edx, dword ptr [rbp-0xc]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37253,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov eax, dword ptr [rip+0x2e8d]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37259,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov r8d, esi"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37262,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov esi, eax"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37264,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"lea rax, ptr [rip+0xe71]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37271,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov rdi, rax"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37274,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov eax, 0x0"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37279,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"call 0x55db4eeb9050"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37284,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov eax, 0x0"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37289,"C_line_num":14,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"leave "}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37290,"C_line_num":14,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"ret "}},{"VACCS_FUNCTION_INV":{"Event_num":1,"Line_num":0,"Path":"__NOCSOURCE__","Inv_line_num":0,"C_inv_path":"__NOCSOURCE__","Func_name":"_start","Address":140064791934592}},{"VACCS_FUNCTION_INV":{"Event_num":3,"Line_num":4,"Path":"","Inv_line_num":0,"C_inv_path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Func_name":"main","Address":140728349461256}},{"VACCS_RETURN_ADDR":{"Event_num":3,"Dynamic_link":140064477568400,"Return_address":1,"Function_name":"main"}},{"VACCS_LINE":{"Event_num":4,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c"}},{"VACCS_REGISTER":{"Event_num":4,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461248}},{"VACCS_REGISTER":{"Event_num":5,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rbp","Value":140728349461248}},{"VACCS_REGISTER":{"Event_num":6,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461232}},{"VACCS_REGISTER":{"Event_num":7,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rflags","Value":518}},{"VACCS_VAR_ACCESS":{"Event_num":8,"Line_num":8,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"main","Address":140728349461240,"Var_type":"\u003cpointer\u003e","Name":"p","Points_to":"","Value":"\u003cNULL\u003e"}},{"VACCS_RETURN_ADDR":{"Event_num":8,"Dynamic_link":140064477568400,"Return_address":1,"Function_name":"main"}},{"VACCS_LINE":{"Event_num":9,"Line_num":10,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c"}},{"VACCS_REGISTER":{"Event_num":9,"Line_num":10,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rax","Value":140728349461236}},{"VACCS_VAR_ACCESS":{"Event_num":10,"Line_num":10,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"main","Address":140728349461240,"Var_type":"\u003cpointer\u003e","Name":"p","Points_to":2,"Value":"0x7ffddf478af4"}},{"VACCS_VAR_ACCESS":{"Event_num":11,"Line_num":11,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"*G*","Address":94400410271768,"Var_type":"int","Name":"a","Points_to":"","Value":"1"}},{"VACCS_VAR_ACCESS":{"Event_num":12,"Line_num":12,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"main","Address":140728349461236,"Var_type":"int","Name":"b","Points_to":"","Value":"2"}},{"VACCS_VAR_ACCESS":{"Event_num":12,"Line_num":12,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"main","Address":140728349461240,"Var_type":"\u003cpointer\u003e","Name":"p","Points_to":2,"Value":"0x7ffddf478af4"}},{"VACCS_LINE":{"Event_num":13,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c"}},{"VACCS_REGISTER":{"Event_num":13,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"esi","Value":3}},{"VACCS_REGISTER":{"Event_num":14,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"ecx","Value":2}},{"VACCS_REGISTER":{"Event_num":15,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"edx","Value":2}},{"VACCS_REGISTER":{"Event_num":16,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"eax","Value":1}},{"VACCS_REGISTER":{"Event_num":17,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"r8d","Value":3}},{"VACCS_REGISTER":{"Event_num":18,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"esi","Value":1}},{"VACCS_REGISTER":{"Event_num":19,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rax","Value":94400410263560}},{"VACCS_REGISTER":{"Event_num":20,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rdi","Value":94400410263560}},{"VACCS_REGISTER":{"Event_num":21,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"eax","Value":0}},{"VACCS_REGISTER":{"Event_num":23,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rip","Value":94400410259536}},{"VACCS_REGISTER":{"Event_num":24,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461224}},{"VACCS_LINE":{"Event_num":26,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c"}},{"VACCS_REGISTER":{"Event_num":26,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rbp","Value":1}},{"VACCS_REGISTER":{"Event_num":27,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461256}},{"VACCS_RETURN":{"Event_num":28}},{"VACCS_REGISTER":{"Event_num":30,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rip","Value":140064477568400}},{"VACCS_REGISTER":{"Event_num":31,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461264}},{"VACCS_RETURN":{"Event_num":35}}];

function generateCCodeMap(VACCS_CCODE) {
  var cCodeMap = {};
  VACCS_CCODE.forEach(element => {
    const cFileName = element.VACCS_CCODE.Path;
    const cLineNum = element.VACCS_CCODE.Line_num;
    if (!cCodeMap[cFileName]) {
      cCodeMap[cFileName] = [];
    }
    cCodeMap[cFileName][cLineNum] = element.VACCS_CCODE;
  });
  return cCodeMap;
}

function getMainFile(VACCS_FUNCTION_INV) {
  for (let i = 0; i < VACCS_FUNCTION_INV.length; i++) {
    if (VACCS_FUNCTION_INV[i].VACCS_FUNCTION_INV.Func_name.trim() === "main") {
      return VACCS_FUNCTION_INV[i].VACCS_FUNCTION_INV.C_inv_path;
    }
  }
  return null;
}

function generateCToASMap(VACCS_ASM) {
  var cToASMMap = {};
  VACCS_ASM.forEach(element => {
    if (!cToASMMap[element.VACCS_ASM.C_file_name])
      cToASMMap[element.VACCS_ASM.C_file_name] = element.VACCS_ASM.Asm_file_name;
  });
  return cToASMMap;
}

function generateASMCodeMap(VACCS_ASM) {
  var asmCodeMap = {};
  VACCS_ASM.forEach(element => {
    const asmFileName = element.VACCS_ASM.Asm_file_name;
    const asmLineNum = element.VACCS_ASM.Asm_line_num;
    if (!asmCodeMap[asmFileName]) {
      asmCodeMap[asmFileName] = [];
    }
    asmCodeMap[asmFileName][asmLineNum] = element.VACCS_ASM.Inst;
  });
  return asmCodeMap;
}

function generateEvents(VACCS_ARCH, VACCS_ASM, VACCS_BINARY, VACCS_CCODE, VACCS_CMD_LINE, VACCS_FILE_CLOSE, VACCS_FILE_OPEN, VACCS_FILE_READ, VACCS_FILE_WRITE, VACCS_FREE, VACCS_FUNCTION_INV, VACCS_LINE, VACCS_MALLOC, VACCS_REGISTER, VACCS_RETURN, VACCS_RETURN_ADDR, VACCS_SD_CLEAR, VACCS_SD_CZERO, VACCS_SD_LOCK, VACCS_SD_SET, VACCS_SD_UNLOCK, VACCS_VAR_ACCESS) {
  var events = [];
  VACCS_FILE_CLOSE.forEach(element => {
    events[element.VACCS_FILE_CLOSE.Event_num] = addToArray(events[element.VACCS_FILE_CLOSE.Event_num], element);
  });
  VACCS_FILE_OPEN.forEach(element => {
    events[element.VACCS_FILE_OPEN.Event_num] = addToArray(events[element.VACCS_FILE_OPEN.Event_num], element);
  });
  VACCS_FILE_READ.forEach(element => {
    events[element.VACCS_FILE_READ.Event_num] = addToArray(events[element.VACCS_FILE_READ.Event_num], element);
  });
  VACCS_FILE_WRITE.forEach(element => {
    events[element.VACCS_FILE_WRITE.Event_num] = addToArray(events[element.VACCS_FILE_WRITE.Event_num], element);
  });
  VACCS_FREE.forEach(element => {
    events[element.VACCS_FREE.Event_num] = addToArray(events[element.VACCS_FREE.Event_num], element);
  });
  VACCS_MALLOC.forEach(element => {
    events[element.VACCS_MALLOC.Event_num] = addToArray(events[element.VACCS_MALLOC.Event_num], element);
  });
  VACCS_REGISTER.forEach(element => {
    events[element.VACCS_REGISTER.Event_num] = addToArray(events[element.VACCS_REGISTER.Event_num], element);
  });
  VACCS_RETURN.forEach(element => {
    events[element.VACCS_RETURN.Event_num] = addToArray(events[element.VACCS_RETURN.Event_num], element);
  });
  VACCS_RETURN_ADDR.forEach(element => {
    events[element.VACCS_RETURN_ADDR.Event_num] = addToArray(events[element.VACCS_RETURN_ADDR.Event_num], element);
  });
  VACCS_SD_CLEAR.forEach(element => {
    events[element.VACCS_SD_CLEAR.Event_num] = addToArray(events[element.VACCS_SD_CLEAR.Event_num], element);
  });
  VACCS_SD_CZERO.forEach(element => {
    events[element.VACCS_SD_CZERO.Event_num] = addToArray(events[element.VACCS_SD_CZERO.Event_num], element);
  });
  VACCS_SD_LOCK.forEach(element => {
    events[element.VACCS_SD_LOCK.Event_num] = addToArray(events[element.VACCS_SD_LOCK.Event_num], element);
  });
  VACCS_SD_SET.forEach(element => {
    events[element.VACCS_SD_SET.Event_num] = addToArray(events[element.VACCS_SD_SET.Event_num], element);
  });
  VACCS_SD_UNLOCK.forEach(element => {
    events[element.VACCS_SD_UNLOCK.Event_num] = addToArray(events[element.VACCS_SD_UNLOCK.Event_num], element);
  });
  VACCS_VAR_ACCESS.forEach(element => {
    events[element.VACCS_VAR_ACCESS.Event_num] = addToArray(events[element.VACCS_VAR_ACCESS.Event_num], element);
  });
  VACCS_FUNCTION_INV.forEach(element => {
    events[element.VACCS_FUNCTION_INV.Event_num] = addToArray(events[element.VACCS_FUNCTION_INV.Event_num], element);
  });
  VACCS_LINE.forEach(element => {
    events[element.VACCS_LINE.Event_num] = addToArray(events[element.VACCS_LINE.Event_num], element);
  });
  return events;
}

function generateSections(VACCS_SECTION) {
  var sections = [];
  VACCS_SECTION.forEach(element => {
    sections.push(parseSectionInfo(element.VACCS_SECTION.Info));
  });
  return sections;
}

function App({ id }) {

  const [analysisData, setAnalysisData] = useState(null);

  useEffect(() => {
    axios.get('/project/analyze/' + id)
      .then(response => {
        setAnalysisData(response.data);
      })
      .catch(error => {
        console.error('There was an error fetching the analysis data!', error);
      });
  }, []);

  const VACCS_ARCH = analysisData.filter(item => item.VACCS_ARCH);  // filter out VACCS_ARCH objects 
  const VACCS_ASM = analysisData.filter(item => item.VACCS_ASM);  // filter out VACCS_ASM objects 
  const VACCS_BINARY = analysisData.filter(item => item.VACCS_BINARY);  // filter out VACCS_BINARY objects
  const VACCS_CCODE = analysisData.filter(item => item.VACCS_CCODE);  // filter out VACCS_CCODE objects 
  const VACCS_CMD_LINE = analysisData.filter(item => item.VACCS_CMD_LINE);  // filter out VACCS_CMD_LINE objects
  const VACCS_FILE_CLOSE = analysisData.filter(item => item.VACCS_FILE_CLOSE);  // filter out VACCS_FILE_CLOSE objects
  const VACCS_FILE_OPEN = analysisData.filter(item => item.VACCS_FILE_OPEN);  // filter out VACCS_FILE_OPEN objects
  const VACCS_FILE_READ = analysisData.filter(item => item.VACCS_FILE_READ);  // filter out VACCS_FILE_READ objects 
  const VACCS_FILE_WRITE = analysisData.filter(item => item.VACCS_FILE_WRITE);  // filter out VACCS_FILE_WRITE objects
  const VACCS_FREE = analysisData.filter(item => item.VACCS_FREE);  // filter out VACCS_FREE objects
  const VACCS_FUNCTION_INV = analysisData.filter(item => item.VACCS_FUNCTION_INV);  // filter out VACCS_FUNCTION_INV objects  
  const VACCS_LINE = analysisData.filter(item => item.VACCS_LINE);  // filter out VACCS_LINE objects  
  const VACCS_MALLOC = analysisData.filter(item => item.VACCS_MALLOC);  // filter out VACCS_MALLOC objects
  const VACCS_REGISTER = analysisData.filter(item => item.VACCS_REGISTER);  // filter out VACCS_REGISTER objects
  const VACCS_RETURN = analysisData.filter(item => item.VACCS_RETURN);  // filter out VACCS_RETURN objects
  const VACCS_RETURN_ADDR = analysisData.filter(item => item.VACCS_RETURN_ADDR);  // filter out VACCS_RETURN_ADDR objects
  const VACCS_SD_CLEAR = analysisData.filter(item => item.VACCS_SD_CLEAR);  // filter out VACCS_SD_CLEAR objects  
  const VACCS_SD_CZERO = analysisData.filter(item => item.VACCS_SD_CZERO);  // filter out VACCS_SD_CZERO objects
  const VACCS_SD_LOCK = analysisData.filter(item => item.VACCS_SD_LOCK);  // filter out VACCS_SD_LOCK objects
  const VACCS_SD_SET = analysisData.filter(item => item.VACCS_SD_SET);  // filter out VACCS_SD_SET objects
  const VACCS_SD_UNLOCK = analysisData.filter(item => item.VACCS_SD_UNLOCK);  // filter out VACCS_SD_UNLOCK objects
  const VACCS_SECTION = analysisData.filter(item => item.VACCS_SECTION);  // filter out VACCS_SECTION objects
  const VACCS_VAR_ACCESS = analysisData.filter(item => item.VACCS_VAR_ACCESS);  // filter out VACCS_VAR_ACCESS objects

  console.log("analysis data = " + analysisData);
  const events = generateEvents(VACCS_ARCH, VACCS_ASM, VACCS_BINARY, VACCS_CCODE, VACCS_CMD_LINE, VACCS_FILE_CLOSE, VACCS_FILE_OPEN, VACCS_FILE_READ, VACCS_FILE_WRITE, VACCS_FREE, VACCS_FUNCTION_INV, VACCS_LINE, VACCS_MALLOC, VACCS_REGISTER, VACCS_RETURN, VACCS_RETURN_ADDR, VACCS_SD_CLEAR, VACCS_SD_CZERO, VACCS_SD_LOCK, VACCS_SD_SET, VACCS_SD_UNLOCK, VACCS_SECTION, VACCS_VAR_ACCESS);
  const sections = generateSections(VACCS_SECTION);
  const asmMap = new ASMSourceMap(VACCS_ASM);
  const cToASMMap = generateCToASMap(VACCS_ASM);
  const cSourceMap = generateCCodeMap(VACCS_CCODE);
  const asmSourceMap = generateASMCodeMap(VACCS_ASM);
  const mainFile = getMainFile(VACCS_FUNCTION_INV);


  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
      <div style={{ flex: '0 0 10%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
        <div style={{ flex: '1', display: 'flex', justifyContent: 'flex-start' }}>
          <Options />
        </div>
        <div style={{ flex: '1', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <div style={{ width: '500px' }}></div>
        </div>
        <div style={{ flex: '1', display: 'flex', justifyContent: 'center' }}>
          <img src={`${process.env.PUBLIC_URL}/splash.png`} alt="Splash" width="200" height="67" />
        </div>
        <div style={{ flex: '1', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <div style={{ width: '500px' }}></div>
        </div>
        <div style={{ flex: '1', display: 'flex', justifyContent: 'flex-end' }}>
          <About />
        </div>
      </div>
      <div style={{ flex: '0 0 100%', width: '98%' }}>
        <ControlPanel cmap={cSourceMap} mainFile={mainFile} ctoasm={cToASMMap} asmsource={asmSourceMap} asmMap={asmMap} sections={sections} events={events} />
      </div>
    </div>
  );
}

export default App;