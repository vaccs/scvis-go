import React, { useState, useEffect } from 'react';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';
import CSource from './CSource';
import {ASMSource, ASMSourceMap} from './ASMSource';
import { SectionTable } from './Sections';
import Registers from './Registers';
import { FunctionInvocationEvent, RegisterEvent, ReturnAddressEvent, ReturnEvent, VariableAccessEvent } from './Event';
import { StackView, StackFrame , RuntimeStack } from './Stack';
import { Globals, GlobalVariables } from './GlobalVariables';

const ControlPanel = ({ cmap, mainFile, ctoasm, asmsource, asmMap, sections, events }) => {
    const initialRegTable = {
        rax: "--",
        rbx: "--",
        rcx: "--",
        rdx: "--",
        rsi: "--",
        rdi: "--",
        rbp: "--",
        rsp: "--",
        r8: "--",
        r9: "--",
        r10: "--",
        r11: "--",
        r12: "--",
        r13: "--",
        r14: "--",
        r15: "--",
        rip: "--",
        eflags: "--"
    };
    const [timeoutId, setTimeoutId] = useState(-1);
    const [eventInfo, setEventInfo] = useState({
        eventNum: 1,
        lineNum: 0,
        path: mainFile,
        isPlaying: false,
        executionRate: 1.0,
        regTable: initialRegTable,
        stack: new RuntimeStack(),
        globals: new Globals() 
    });

    useEffect(() => {
        if (eventInfo.isPlaying) {
            if (eventInfo.eventNum > 0 && eventInfo.eventNum < events.length) {
                const event = events[eventInfo.eventNum];
                var info = handleEvent(event);
                const timeoutId = setTimeout(() => {
                    updateEvent(info)
                }, 1000 / eventInfo.executionRate);
                setTimeoutId(timeoutId)
            } else {
                setEventInfo({
                    eventNum: eventInfo.eventNum,
                    lineNum: eventInfo.lineNum,
                    path: eventInfo.path,
                    isPlaying: false,
                    executionRate: eventInfo.executionRate,
                    regTable: eventInfo.regTable,
                    stack: eventInfo.stack,
                    globals: eventInfo.globals
                });
            }
        }
    }, [eventInfo]);

    const handleEvent = (event) => {
        var info = {
            eventNum: eventInfo.eventNum,
            path: eventInfo.path,
            isPlaying: eventInfo.isPlaying,
            executionRate: eventInfo.executionRate,
            regTable: eventInfo.regTable,
            stack: eventInfo.stack,
            globals: eventInfo.globals
        };
        if (event) {
            for (let i = 0; i < event.length; i++) {
                if (event[i] instanceof RegisterEvent) {
                    info.regTable[event[i].registerName] = event[i].value;
                } else if (event[i] instanceof FunctionInvocationEvent) {
                    info.stack.push({StackFrame: new StackFrame(event[i].functionName, event[i].address)});
                } else if (event[i] instanceof ReturnAddressEvent) {
                    info.stack.updateFrameLinks(event[i].functionName,event[i].address, event[i].dynamicLink);
                }
                else if (event[i] instanceof ReturnEvent) {
                    info.stack.pop();
                } else if (event[i] instanceof VariableAccessEvent) {
                    if (event[i].scope === "*G*") {
                        info.globals.updateVariable(event[i]);
                    } else {
                        info.stack.updateStackVariable(event[i]);
                    }
                }
                if (event[i].lineNum)
                    info.lineNum = event[i].lineNum;
                if (event[i].path)
                    info.path = event[i].path;  
            }
        }
        info.eventNum++;
        return info;
    };

    const clearEvent = () => {
        if (timeoutId !== -1) {
            clearTimeout(timeoutId);
            setTimeoutId(-1);
        }
    }

    const updateEvent = (info) => {
        clearEvent();
        setEventInfo({
            eventNum: info.eventNum,
            lineNum: info.lineNum,
            path: info.path,
            isPlaying: info.isPlaying,
            executionRate: info.executionRate,
            regTable: info.regTable,
            stack: info.stack,
            globals: info.globals
        });
    };

    const handleStart = () => {
        clearEvent();
        setEventInfo({ eventNum: 1, lineNum: 0, path: mainFile, isPlaying: false, executionRate: eventInfo.executionRate, 
            regTable: initialRegTable, stack: new RuntimeStack(), globals: new Globals() });
    };


    const handlePlayPause = () => {
        clearEvent();
        setEventInfo({
            eventNum: eventInfo.eventNum,
            lineNum: eventInfo.lineNum,
            path: eventInfo.path,
            isPlaying: !eventInfo.isPlaying,
            executionRate: eventInfo.executionRate,
            regTable: eventInfo.regTable,
            stack: eventInfo.stack,
            globals: eventInfo.globals
        });
    };

    const handleEnd = () => {
        clearEvent();
        setEventInfo({ eventNum: events.length, lineNum: eventInfo.lineNum, path: mainFile, isPlaying: false, executionRate: eventInfo.executionRate, 
            regTable: eventInfo.regTable, stack: new RuntimeStack(), globals: new Globals() });
    };

    const handleBack = () => {
        clearEvent();
    };

    const handleNext = () => {
        clearEvent();
    };

    const handleSliderChange = (event) => {
        setEventInfo({ eventNum: eventInfo.eventNum, lineNum: eventInfo.lineNum, path: eventInfo.path, isPlaying: eventInfo.isPlaying, 
            executionRate: event.target.value,
            regTable: eventInfo.regTable,
            stack: eventInfo.stack,
            globals: eventInfo.globals});
    };

    return (
        <div>
            <div style={{ flex: '0 0 70%', width: '100%' }}>
                <div style={{ display: 'flex', height: '70vh' }}>
                    <div style={{ flex: 1, borderRight: '1px solid #ccc' }}>
                        <Tabs>
                            <TabList>
                                <Tab>Global Variables</Tab>
                                <Tab>Stack</Tab>
                                <Tab>Registers</Tab>
                                <Tab>Sections</Tab>
                                <Tab style={{ pointerEvents: 'none', opacity: 0.5 }}>Heap</Tab>
                                <Tab style={{ pointerEvents: 'none', opacity: 0.5 }}>Sensitive Data</Tab>
                            </TabList>

                            <TabPanel>
                                <GlobalVariables globals={eventInfo.globals}/>
                            </TabPanel>
                            <TabPanel>
                                <StackView stack={eventInfo.stack} />
                            </TabPanel>
                            <TabPanel>
                                <Registers registers={eventInfo.regTable} />
                            </TabPanel>
                            <TabPanel>
                                <SectionTable sections={sections} />
                            </TabPanel>
                            <TabPanel>
                                
                            </TabPanel>
                            <TabPanel>
                                
                            </TabPanel>
                        </Tabs>
                    </div>
                    <div style={{ flex: 1 }}>
                        <Tabs>
                            <TabList>
                                <Tab>C Source</Tab>
                                <Tab>ASM Source</Tab>
                                <Tab>Output</Tab>
                                <Tab style={{ pointerEvents: 'none', opacity: 0.5 }}>File Operations</Tab>
                                <Tab style={{ pointerEvents: 'none', opacity: 0.5 }}>Call Graph</Tab>
                            </TabList>

                            <TabPanel>
                                <CSource cmap={cmap} path={eventInfo.path} line={eventInfo.lineNum} />
                            </TabPanel>
                            <TabPanel>
                                <ASMSource asmsource={asmsource} ctoasm={ctoasm} asmMap={asmMap} cfile={eventInfo.path} line={eventInfo.lineNum} />
                            </TabPanel>
                            <TabPanel>
                                
                            </TabPanel>
                            <TabPanel>
                                <h2>File Operations Content</h2>
                            </TabPanel>
                            <TabPanel>
                                <h2>Call Graph Content</h2>
                            </TabPanel>
                        </Tabs>
                    </div>
                </div>
            </div>
            <div style={{ flex: '0 0 20%', width: '100%' }}>
                <div className="control-panel-container" style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', height: '20vh' }}>
                    <div className="slider" style={{ marginBottom: '10px', display: 'flex', alignItems: 'center' }}>
                        <span style={{ marginRight: '10px' }}>0.5x</span>
                        <input
                            type="range"
                            min="0.5"
                            max="2"
                            step="0.1"
                            value={eventInfo.executionRate}
                            onChange={handleSliderChange}
                            style={{ flexGrow: 1 }}
                        />
                        <span style={{ marginLeft: '10px' }}>2x</span>
                    </div>
                    <div className="control-panel">
                        <div className="buttons" style={{ display: 'flex', justifyContent: 'center' }}>
                            <button onClick={handleStart}><img src="/playctl_start.png" alt="Skip to Start" /></button>
                            <button onClick={handleBack}> <img src="/playctl_back.png" alt="Back" /> </button>
                            <button onClick={handlePlayPause}>{eventInfo.isPlaying ? <img src="/pause.png" alt="Pause by Hilmy Abiyyu A." width="64" height="64" /> : <img src="/playctl_play.png" alt="Play" />}</button>
                            <button onClick={handleNext}><img src="/playctl_next.png" alt="Next" /></button>
                            <button onClick={handleEnd}><img src="/playctl_end.png" alt="Skip to End" /></button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ControlPanel;