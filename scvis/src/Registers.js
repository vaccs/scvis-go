import React from 'react';
import PropTypes from 'prop-types';

const Registers = ({ registers }) => {
    return (
        <div style={{display: 'flex', height: '65vh', overflowY: 'scroll' }}>
            <table>
                <thead>
                    <tr>
                        <th>Register Name</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    {Object.entries(registers).map(([name, value], index) => (
                        <tr key={name} style={{ backgroundColor: index % 2 === 0 ? 'brown' : 'white', color: index % 2 === 0 ? 'white' : 'black' }}>
                            <td>{name}</td>
                            <td>{value}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

Registers.propTypes = {
    registers: PropTypes.object.isRequired,
};

export default Registers;