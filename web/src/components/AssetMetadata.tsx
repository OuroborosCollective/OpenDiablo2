import React from 'react';

interface AssetMetadataItem {
  id: string;
  name: string;
  size: string;
  type: string;
  flags: string;
  compression: string;
}

const mockAssets: AssetMetadataItem[] = [
  { id: '1', name: 'data\\global\\ui\\Loading\\BarBackground.dc6', size: '12.4 KB', type: 'DC6', flags: '0x01000000', compression: 'ZLib' },
  { id: '2', name: 'data\\global\\excel\\chars.bin', size: '4.2 KB', type: 'BIN', flags: '0x00000200', compression: 'None' },
  { id: '3', name: 'data\\local\\ui\\eng\\inventory\\invslot.dc6', size: '8.1 KB', type: 'DC6', flags: '0x01000000', compression: 'ZLib' },
  { id: '4', name: 'd2data.mpq', size: '412 MB', type: 'MPQ', flags: 'N/A', compression: 'Multi' },
  { id: '5', name: 'd2exp.mpq', size: '185 MB', type: 'MPQ', flags: 'N/A', compression: 'Multi' },
];

const AssetMetadata: React.FC = () => {
  return (
    <div className="flex flex-col w-full md:w-96 bg-gray-900 text-gray-100 h-full border-l border-gray-700 overflow-y-auto">
      <div className="p-4 border-b border-gray-700">
        <h2 className="text-xl font-bold text-orange-500 uppercase tracking-wider">Asset Metadata</h2>
        <p className="text-[10px] text-gray-500 mt-1 uppercase">Axiomatic MPQ Inspection</p>
      </div>
      <div className="flex-1">
        {mockAssets.map((asset) => (
          <div key={asset.id} className="p-4 border-b border-gray-800 hover:bg-gray-800 transition-colors cursor-pointer group">
            <div className="text-sm font-mono text-gray-300 break-all group-hover:text-white">
              {asset.name}
            </div>
            <div className="grid grid-cols-2 gap-2 mt-3 text-[10px] font-mono">
              <div className="flex flex-col">
                <span className="text-gray-500 uppercase">Type</span>
                <span className="text-orange-400">{asset.type}</span>
              </div>
              <div className="flex flex-col">
                <span className="text-gray-500 uppercase">Size</span>
                <span className="text-blue-400">{asset.size}</span>
              </div>
              <div className="flex flex-col">
                <span className="text-gray-500 uppercase">Flags</span>
                <span className="text-green-400">{asset.flags}</span>
              </div>
              <div className="flex flex-col">
                <span className="text-gray-500 uppercase">Comp</span>
                <span className="text-purple-400">{asset.compression}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
      <div className="p-4 bg-black/50 border-t border-gray-800">
        <div className="flex justify-between items-center text-[10px] text-gray-600 italic">
          <span>* Recursive BaalAal Logic:</span>
          <span className="font-mono text-orange-900">{new Date().getTime().toString(16)}</span>
        </div>
      </div>
    </div>
  );
};

export default AssetMetadata;
