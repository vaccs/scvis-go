class ASMInfo {
    constructor(asmfile, rangeStart, rangeEnd) {
        this.asmfile = asmfile;
        this.rangeStart = rangeStart;
        this.rangeEnd = rangeEnd;

    }
}

class ASMSourceMap {
    constructor(asmdata) {
        this.asmdata = asmdata;
        this.map = new Map();
        this.buildMap();
    }

    add(cFileName, cLineNumber, asmInfo) {
        if (!this.map.has(cFileName)) {
            this.map.set(cFileName, []);
        }
        this.map.get(cFileName)[cLineNumber] = asmInfo;
    }

    get(cFileName, cLineNumber) {
        if (this.map.has(cFileName)) {
            return this.map.get(cFileName)[cLineNumber];
        }
        return null;
    }

    buildMap() {
        this.asmdata.forEach(item => {
            const cFileName = item.VACCS_ASM.C_file_name;
            const cLineNumber = item.VACCS_ASM.C_line_num;  
            var info = this.get(cFileName, cLineNumber);
            if (!info) {
                info = new ASMInfo(item.VACCS_ASM.Asm_file_name, item.VACCS_ASM.Asm_line_num, item.VACCS_ASM.Asm_line_num);
            } else {
                if (item.VACCS_ASM.Asm_line_num < info.rangeStart) {
                    info.rangeStart = item.VACCS_ASM.Asm_line_num;
                } else if (item.VACCS_ASM.Asm_line_num > info.rangeEnd) {
                    info.rangeEnd = item.VACCS_ASM.Asm_line_num;
                }
            }
            this.add(cFileName, cLineNumber, info);
        });
    }
}

const ASMSource = ({asmsource, ctoasm, asmMap, cfile, line}) => {
    const asmfile = ctoasm[cfile];
    var asmInfo = asmMap.get(cfile, line);
    const source = asmsource[asmfile];

    if (!asmInfo) {
        asmInfo = new ASMInfo('', 0);
    }

    var line = 1;

    return (
        <div>
            {source && source.map((item, index) => (
                <div key={index} style={{ 
                        whiteSpace: 'pre-wrap', 
                        backgroundColor: ((index+1) >= asmInfo.rangeStart && (index+1) <= asmInfo.rangeEnd) ? 'brown' : 'inherit',
                        color: ((index+1) >= asmInfo.rangeStart && (index+1) <= asmInfo.rangeEnd) ? 'white' : 'inherit',
                    }}>
                    {line++ +'\t' + item}
                </div>
            ))}
        </div>
    );
};

export { ASMSource, ASMSourceMap, ASMInfo };