import React from 'react';

class VariableInfo {
    constructor(name, value, address, type, point_to = null) {
        this.name = name;
        this.value = value;
        this.address = address;
        this.type = type;
        this.point_to = point_to;
    }

    updateVar(value, pointsTo) {
        this.value = value;
        this.point_to = pointsTo;
    }
}

class StackFrame {
    constructor(functionName, address) {
        this.functionName = functionName;
        this.address = address;
        this.returnAddress = "<unknown>";
        this.dynamicLink = "<unknown>";
        this.variables = [];
    }

    addVariable(variable) {
        let insertIndex = this.variables.findIndex(v => v.address > variable.address);
        if (insertIndex === -1) {
            insertIndex = this.variables.length;
        }
        this.variables.splice(insertIndex, 0, variable);
    }

    findVariable(variableName) {
        return this.variables.find(variable => variable.name === variableName);
    }

    isScope(scope) {
        return this.functionName === scope;
    }

    updateFrameVariables(variableAccessEvent) {
        var varInfo = this.findVariable(variableAccessEvent.variableName);
        if (varInfo) {
            varInfo.updateVar(variableAccessEvent.value, variableAccessEvent.pointsTo);
        } else
            this.addVariable(new VariableInfo(variableAccessEvent.variableName, variableAccessEvent.value, variableAccessEvent.address, variableAccessEvent.varType, variableAccessEvent.pointsTo));
    }

    updateLinks(returnAddress, dynamicLink) {
        this.returnAddress = returnAddress;
        this.dynamicLink = dynamicLink;
    }
}

class RuntimeStack {
    constructor() {
        this.stack = [];
    }

    push(stackFrame) {
        this.stack.unshift(stackFrame);
    }

    pop() {
        return this.stack.shift();
    }

    replace(stackFrame) {
        const index = this.stack.findIndex(frame => frame.StackFrame.isScope(stackFrame.functionName));
        if (index !== -1) {
            this.stack[index] = stackFrame;
        }
    }

    getStack() {
        return this.stack;
    }

    findStackFrameByScope(scope) {
        return this.stack.find(frame => frame.StackFrame.isScope(scope));
    }

    updateFrameLinks(name, address, dynamicLink) {
        var frame = this.findStackFrameByScope(name);
        if (frame) {
            frame.StackFrame.updateLinks(address, dynamicLink);
        }
    }

    updateStackVariable(variableAccessEvent) {
        var frame = this.findStackFrameByScope(variableAccessEvent.scope);

        if (frame) {
            frame.StackFrame.updateFrameVariables(variableAccessEvent);
        }
    }
}




const StackView = ({ stack }) => {
    const getRowStyle = (index) => {
        switch (index % 3) {
            case 0:
                return { backgroundColor: 'brown', color: 'white' };
            case 1:
                return { backgroundColor: 'green', color: 'white' };
            case 2:
                return { backgroundColor: 'blue', color: 'white' };
            default:
                return {};
        }
    };
    return (
        <div style={{ display: 'flex', maxHeight: '65vh', overflowY: 'scroll' }}>
            {stack !== undefined && stack.stack.length > 0 && (
                <table style={{ borderCollapse: 'collapse', width: '98%' }}>
                    <thead>
                        <tr>
                            <th style={{ border: '1px solid black', padding: '8px' }}>Name</th>
                            <th style={{ border: '1px solid black', padding: '8px' }}>Value</th>
                            <th style={{ border: '1px solid black', padding: '8px' }}>Address</th>
                            <th style={{ border: '1px solid black', padding: '8px' }}>Type</th>
                            <th style={{ border: '1px solid black', padding: '8px' }}>At Target Address</th>
                        </tr>
                    </thead>
                    <tbody>
                        {stack.getStack().map((frame, index) => (
                            <>
                                <tr key={index} style={{ ...getRowStyle(index), border: '1px solid black' }}>
                                    <td colSpan="5" style={{ border: '1px solid black', padding: '8px' }}>Frame: {frame.StackFrame.functionName} </td>
                                </tr>
                                {frame.StackFrame.variables && frame.StackFrame.variables.map((variable) => (
                                    <tr key={`${index}-${variable.name}`} style={{ backgroundColor: 'white', color: 'black', border: '1px solid black', padding: '8px' }}>
                                        <td style={{ border: '1px solid black', padding: '8px' }}>{variable.name}</td>
                                        <td style={{ border: '1px solid black', padding: '8px' }}>{variable.value}</td>
                                        <td style={{ border: '1px solid black', padding: '8px' }}>{"0x" + variable.address.toString(16)}</td>
                                        <td style={{ border: '1px solid black', padding: '8px' }}>{variable.type}</td>
                                        <td style={{ border: '1px solid black', padding: '8px' }}>{variable.type === "<pointer>" ? variable.point_to : ""}</td>
                                    </tr>
                                ))}
                                <tr key={`${index}-dynamic`} style={{ backgroundColor: 'white', color: 'black', border: '1px solid black' }}>
                                    <td style={{ border: '1px solid black', padding: '8px' }}>Dynamic Link: </td>
                                    <td style={{ border: '1px solid black', padding: '8px' }}>{isNaN(frame.StackFrame.dynamicLink) ? frame.StackFrame.dynamicLink : "0x" + frame.StackFrame.dynamicLink.toString(16)}</td>
                                    <td colSpan="3" style={{ border: '1px solid black', padding: '8px' }}></td>
                                </tr>
                                <tr key={`${index}-return`} style={{ backgroundColor: 'white', color: 'black', border: '1px solid black' }}>
                                    <td style={{ border: '1px solid black', padding: '8px' }}>Return Address: </td>
                                    <td style={{ border: '1px solid black', padding: '8px' }}>{isNaN(frame.StackFrame.returnAddress) ? frame.StackFrame.returnAddress : "0x" + frame.StackFrame.returnAddress.toString(16)}</td>
                                    <td colSpan="3" style={{ border: '1px solid black', padding: '8px' }}></td>
                                </tr>
                            </>
                        ))}
                    </tbody>
                </table>)}
        </div>
    );
};

export { StackView, StackFrame, RuntimeStack, VariableInfo };