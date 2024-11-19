import { VariableInfo } from './Stack';

class Globals {
    constructor() {
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
        return scope === "*G*";
    }

    updateVariable(variableAccessEvent) {
        var varInfo = this.findVariable(variableAccessEvent.variableName);
        if (varInfo) {
            varInfo.updateVar(variableAccessEvent.value, variableAccessEvent.pointsTo);
        } else
            this.addVariable(new VariableInfo(variableAccessEvent.variableName, variableAccessEvent.value, variableAccessEvent.address, variableAccessEvent.varType, variableAccessEvent.pointsTo));
    }
}

const GlobalVariables = ({ globals }) => {
    return (
        <div style={{ display: 'flex', maxHeight: '65vh', overflowY: 'scroll' }}>
            {globals !== undefined && globals.variables.length > 0 && (
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
                        {globals.variables.map((variable) => (
                            <tr key={`${variable.name}`} style={{ backgroundColor: 'white', color: 'black' }}>
                                <td style={{ border: '1px solid black', padding: '8px' }}>{variable.name}</td>
                                <td style={{ border: '1px solid black', padding: '8px' }}>{variable.value}</td>
                                <td style={{ border: '1px solid black', padding: '8px' }}>{"0x" + variable.address.toString(16)}</td>
                                <td style={{ border: '1px solid black', padding: '8px' }}>{variable.type}</td>
                                <td style={{ border: '1px solid black', padding: '8px' }}>{variable.type === "<pointer>" ? variable.point_to : ""}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>)}
        </div>
    );
};

export { GlobalVariables, Globals };