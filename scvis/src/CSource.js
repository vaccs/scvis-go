const CSource = ({cmap, path, line}) => {
    return (
        <div>
            {cmap[path] && cmap[path].map((item, index) => (
                <div 
                    key={index} 
                    style={{ 
                        whiteSpace: 'pre-wrap', 
                        backgroundColor: (index+1) === line ? 'brown' : 'inherit',
                        color: (index+1) === line ? 'white' : 'inherit',
                    }}
                >
                    {index + 1 + '\t' + item.Source}
                </div>
            ))}
        </div>
    );
};

export default CSource;