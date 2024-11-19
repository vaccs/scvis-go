import React, { useState, useEffect } from 'react';
import Popup from 'reactjs-popup';
import 'reactjs-popup/dist/index.css';

const Options = () => {
    const [detailMode, setDetailMode] = useState(localStorage.getItem('detailMode') || 'Novice');
    const [viewElements, setViewElements] = useState(JSON.parse(localStorage.getItem('viewElements')) || {
        simplifiedPAS: false,
        hiddenFunctions: false,
        cpuRegisters: false,
        allProgramSections: false,
        showText: false,
        assemblyCode: false,
        programOutput: false,
        offsets: false,
    });

    useEffect(() => {
        localStorage.setItem('detailMode', detailMode);
    }, [detailMode]);

    useEffect(() => {
        localStorage.setItem('viewElements', JSON.stringify(viewElements));
    }, [viewElements]);

    const handleRadioChange = (event) => {
        setDetailMode(event.target.value);
    };

    const handleCheckboxChange = (event) => {
        setViewElements({
            ...viewElements,
            [event.target.name]: event.target.checked,
        });
    };

    return (
        <Popup trigger={<button className="button"> Options </button>} modal>
            {close => (
                <div className="popup-content">
                    <div className="column">
                        <h3>Detail Mode</h3>
                        <label>
                            <input
                                type="radio"
                                value="Novice"
                                checked={detailMode === 'Novice'}
                                onChange={handleRadioChange}
                            />
                            Novice
                        </label>
                        <label>
                            <input
                                type="radio"
                                value="Intermediate"
                                checked={detailMode === 'Intermediate'}
                                onChange={handleRadioChange}
                            />
                            Intermediate
                        </label>
                        <label>
                            <input
                                type="radio"
                                value="Advanced"
                                checked={detailMode === 'Advanced'}
                                onChange={handleRadioChange}
                            />
                            Advanced
                        </label>
                        <label>
                            <input
                                type="radio"
                                value="Expert"
                                checked={detailMode === 'Expert'}
                                onChange={handleRadioChange}
                            />
                            Expert
                        </label>
                    </div>
                    <div className="column">
                        <h3>View Elements</h3>
                        <label>
                            <input
                                type="checkbox"
                                name="simplifiedPAS"
                                checked={viewElements.simplifiedPAS}
                                onChange={handleCheckboxChange}
                            />
                            Simplified PAS
                        </label>
                        <label>
                            <input
                                type="checkbox"
                                name="hiddenFunctions"
                                checked={viewElements.hiddenFunctions}
                                onChange={handleCheckboxChange}
                            />
                            Hidden Functions
                        </label>
                        <label>
                            <input
                                type="checkbox"
                                name="cpuRegisters"
                                checked={viewElements.cpuRegisters}
                                onChange={handleCheckboxChange}
                            />
                            CPU Registers
                        </label>
                        <label>
                            <input
                                type="checkbox"
                                name="allProgramSections"
                                checked={viewElements.allProgramSections}
                                onChange={handleCheckboxChange}
                            />
                            All Program Sections
                        </label>
                        <label>
                            <input
                                type="checkbox"
                                name="showText"
                                checked={viewElements.showText}
                                onChange={handleCheckboxChange}
                            />
                            Show .text
                        </label>
                        <label>
                            <input
                                type="checkbox"
                                name="assemblyCode"
                                checked={viewElements.assemblyCode}
                                onChange={handleCheckboxChange}
                            />
                            ASM Code
                        </label>
                        <label>
                            <input
                                type="checkbox"
                                name="programOutput"
                                checked={viewElements.programOutput}
                                onChange={handleCheckboxChange}
                            />
                            Program Output
                        </label>
                        <label>
                            <input
                                type="checkbox"
                                name="offsets"
                                checked={viewElements.offsets}
                                onChange={handleCheckboxChange}
                            />
                            Offsets
                        </label>
                    </div>
                    <div className="popup-buttons">
                        <button className="button" onClick={() => close()}>Done</button>
                    </div>
                </div>
            )}
        </Popup>
    );
};

export default Options;