import React from 'react';
import Popup from 'reactjs-popup';
import 'reactjs-popup/dist/index.css';

const About = () => {
    const text = `Authors: James Walker and Steve Carr\n\nThis work has been supported by the National Science Foundation under grants DUE-1245310, DGE-1522883, DGE-1523017, IIS-1456763, and IIS-1455886.`;

    return (
        <div>
            <Popup trigger={<button className="button"> About </button>} modal>
                <div>
                    {text.split('\n').map((line, index) => (
                        <p key={index}>{line}</p>
                    ))}
                </div>
            </Popup>
        </div>
    );
};

export default About;