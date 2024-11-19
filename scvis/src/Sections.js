import React from 'react';

class SectionInfo {
    constructor(type, startAddress, size, name) {
        this.type = type;
        this.startAddress = startAddress;
        this.size = size;
        this.name = name;
    }
}

function parseSectionInfo(sectionString) {
    const [_, type, startAddress, size, name] = sectionString.split('|');
    return new SectionInfo(parseInt(type,10), startAddress, parseInt(size, 10), name);
}

function SectionTable({ sections }) {
    return (
        <div style={{display: 'flex', height: '65vh', overflowY: 'scroll' }}>
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Start Address</th>
                        <th>Type</th>
                        <th>Size</th>
                    </tr>
                </thead>
                <tbody>
                    {sections.slice(1).map((section, index) => (
                        <tr key={index} style={{ backgroundColor: index % 2 === 0 ? 'brown' : 'white', color: index % 2 === 0 ? 'white' : 'black' }}>
                            <td>{section.name}</td>
                            <td>{section.startAddress}</td>
                            <td>{section.type}</td>
                            <td>{section.size}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
export { SectionTable, SectionInfo, parseSectionInfo };