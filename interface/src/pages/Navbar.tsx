import React, { useState, useEffect } from 'react';
import { ConfirmStartOver, Status } from '../js/wailsjs/go/main/App';
import { model } from '../js/wailsjs/go/models';
import { useNavigate } from 'react-router-dom';
interface NavbarProps {
  showBackIcon?: boolean;
  showForwardIcon?: boolean;
  showSettingIcon?: boolean;
  showStartOver?: boolean;
  onBack?: () => void;
  onForward?: () => void;
  onSetting?: () => void;
  onStartOver?: () => void;
  disableBack?: boolean;
  disableForward?: boolean;
}

const Navbar: React.FC<NavbarProps> = ({
  showBackIcon = false,
  showForwardIcon = false,
  showSettingIcon = true,
  showStartOver = false,
  onBack,
  onForward,
  onSetting,
  onStartOver,
  disableBack = false,
  disableForward = false
}) => {
  const [status, setStatus] = useState<model.StatusResponse | null>(null);
  const [showTooltip, setShowTooltip] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const response = await Status();
        setStatus(response);
      } catch (error) {
        console.error('Error fetching status:', error);
        setStatus(null);
      }
    };

    fetchStatus();
    const intervalId = setInterval(fetchStatus, 30000); // Fetch every 30 seconds

    return () => clearInterval(intervalId);
  }, []);

  const getStatusColor = () => {
    if (!status) return 'bg-red-500';
    if (status.server && status.api === true && status.key) return 'bg-green-500';
    if (status.api === false) return 'bg-orange-500';
    return 'bg-red-500';
  };

  const capitalizeFirstLetter = (string: string) => {
    return string.charAt(0).toUpperCase() + string.slice(1);
  };

  if (!showBackIcon && !showForwardIcon && !showStartOver && !status) {
    return null;
  }

  const handleBack = () => {
    if (!disableBack && onBack) {
      onBack();
    }
  }

  const handleForward = () => {
    if (!disableForward && onForward) {
      onForward();
    }
  }

  const handleStartOver = async () => {
    if (onStartOver) {
      const isConfirmed = await ConfirmStartOver();
      if (isConfirmed === "Yes" || isConfirmed === "Ok" || isConfirmed === "Continue") {
        onStartOver();
      }
    }
  };

  const handleSettings = () => {
    if (onSetting) {
      onSetting();
    }

    navigate('/setting');
  }

  return (
    <nav className="bg-dark-surface p-4 flex justify-between items-center relative">
      <div className="flex items-center">
        {showBackIcon && (
          <button
            onClick={handleBack}
            className={`mr-4 ${disableBack ? 'text-gray-500 cursor-not-allowed' : 'text-white hover:text-gray-300'}`}
            disabled={disableBack}
          >
            &#8592; {/* Left arrow */}
          </button>
        )}
        {showForwardIcon && (
          <button
            onClick={handleForward}
            className={`${disableForward ? 'text-gray-500 cursor-not-allowed' : 'text-white hover:text-gray-300'}`}
            disabled={disableForward}
          >
            &#8594; {/* Right arrow */}
          </button>
        )}
      </div>
      <div className="flex items-center">
        <div 
          className={`w-4 h-4 rounded-full mr-4 ${getStatusColor()}`}
          onMouseEnter={() => setShowTooltip(true)}
          onMouseLeave={() => setShowTooltip(false)}
        >
          {showTooltip && (
            <div className="absolute top-full right-2 mt-2 p-2 bg-white text-black rounded shadow-lg z-10">
              <p>Database: {status?.server ? 'Up' : 'Down'}</p>
              <p>API: {status?.api === null ? 'Unknown' : (status?.api ? 'Up' : 'Down')}</p>
              <p>Status: {status?.apiStatus ? capitalizeFirstLetter(status?.apiStatus) : 'Unknown'}</p>
              <p>Authorized: {status?.key ? 'Yes' : 'No'}</p>
            </div>
          )}
        </div>
        {showSettingIcon && (<button 
          onClick={handleSettings}
          className="mr-4 text-white hover:text-gray-300"
        >
          &#x26ED; {/* Gear icon */}
        </button>)}
        {showStartOver && (
          <button 
            onClick={handleStartOver} 
            className="text-white hover:text-gray-300"
          >
            &#x2715; {/* Multiplication symbol */}
          </button>
        )}
      </div>
    </nav>
  );
};

export default Navbar;